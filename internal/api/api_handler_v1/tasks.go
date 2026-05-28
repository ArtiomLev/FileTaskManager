package api_handler_v1

import (
	"LaserTaskSystem/internal/api/api_v1"
	"LaserTaskSystem/internal/task"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-faster/errors"
)

var forbiddenNames = map[string]bool{
	"":          true,
	".":         true,
	"..":        true,
	"meta.yaml": true,
}

func isValidName(name string) bool {
	if forbiddenNames[name] {
		return false
	}
	if strings.ContainsAny(name, "/\\") {
		return false
	}
	return true
}

func (h *TaskHandler) TasksGet(ctx context.Context, params api_v1.TasksGetParams) (api_v1.TasksGetRes, error) {
	// Get manager
	manager, ok := h.managers[params.ManagerName]
	if !ok {
		return &api_v1.TasksGetNotFound{
			Code:    api_v1.ErrorResponseCodeNOTFOUND,
			Message: "Manager not found",
			Status:  http.StatusNotFound,
		}, nil
	}

	// Get parameters
	contract, hasContract := params.Contract.Get()
	status, hasStatus := params.Status.Get()

	//Check parameters
	if hasContract && !isValidName(contract) {
		return &api_v1.TasksGetBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid contract!",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if hasStatus && status != api_v1.TasksGetStatusActive && status != api_v1.TasksGetStatusCompleted {
		return &api_v1.TasksGetBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid status!",
			Status:  http.StatusBadRequest,
		}, nil
	}

	var tasks []task.Task
	var err error

	switch {
	case hasContract && hasStatus:
		if status == api_v1.TasksGetStatusActive {
			tasks, err = manager.ListActiveContract(contract)
		} else {
			tasks, err = manager.ListCompletedContract(contract)
		}
	case hasContract:
		tasks, err = manager.ListAllContract(contract)
	case hasStatus:
		if status == api_v1.TasksGetStatusActive {
			tasks, err = manager.ListActive()
		} else {
			tasks, err = manager.ListCompleted()
		}
	default:
		tasks, err = manager.ListAll()
	}
	if err != nil {
		return nil, err
	}

	result := make(api_v1.Tasks, len(tasks))
	for i, tsk := range tasks {
		result[i] = api_v1.Task{
			Contract:  tsk.ContractNumber,
			Name:      tsk.Name,
			Completed: tsk.Completed,
			Updated:   tsk.Updated,
		}
	}
	return &result, nil
}

func (h *TaskHandler) TaskGet(ctx context.Context, params api_v1.TaskGetParams) (api_v1.TaskGetRes, error) {
	manager, ok := h.managers[params.ManagerName]
	if !ok {
		return &api_v1.TaskGetNotFound{
			Code:    api_v1.ErrorResponseCodeNOTFOUND,
			Message: "Manager not found",
			Status:  http.StatusNotFound,
		}, nil
	}

	// Get parameters
	contract := params.Contract
	name := params.Name

	//Check parameters
	if !isValidName(contract) {
		return &api_v1.TaskGetBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid contract!",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if !isValidName(name) {
		return &api_v1.TaskGetBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid name!",
			Status:  http.StatusBadRequest,
		}, nil
	}

	// Get task
	tsk, err := manager.GetTask(contract, name)
	if err != nil {
		if errors.Is(err, task.ErrTaskNotExists) {
			return &api_v1.TaskGetNotFound{
				Code:    api_v1.ErrorResponseCodeNOTFOUND,
				Message: "Task not found!",
				Status:  http.StatusNotFound,
			}, nil
		}
		return nil, err
	}

	// Get task files
	files, err := tsk.ListFiles()
	if err != nil {
		return nil, err
	}

	// Make response
	response := api_v1.TaskWithFiles{
		Name:      tsk.Name,
		Contract:  tsk.ContractNumber,
		Completed: tsk.Completed,
		Updated:   tsk.Updated,
		Files:     files,
	}

	return &response, nil
}

func (h *TaskHandler) TaskPost(ctx context.Context, req *api_v1.TaskCreateMultipart, params api_v1.TaskPostParams) (api_v1.TaskPostRes, error) {
	manager, ok := h.managers[params.ManagerName]
	if !ok {
		return &api_v1.TaskPostNotFound{
			Code:    api_v1.ErrorResponseCodeNOTFOUND,
			Message: "Manager not found",
			Status:  http.StatusNotFound,
		}, nil
	}

	// Get task parameters
	name := req.Name
	contract := req.Contract
	files := req.Files

	// Check parameters
	if !isValidName(contract) {
		return &api_v1.TaskPostBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid contract!",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if !isValidName(name) {
		return &api_v1.TaskPostBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid name!",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if len(files) == 0 {
		return &api_v1.TaskPostBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Must have at least one file!",
			Status:  http.StatusBadRequest,
		}, nil
	}

	// Generate file inputs array
	var fileInputs []task.FileInput
	for _, f := range files {
		fileInputs = append(fileInputs, task.FileInput{
			Name:   f.Name,
			Reader: f.File,
		})
	}

	// Create task
	tsk, err := manager.CreateTask(contract, name, fileInputs)
	if err != nil {
		if errors.Is(err, task.ErrTaskExists) {
			return &api_v1.TaskPostConflict{
				Code:    api_v1.ErrorResponseCodeCONFLICT,
				Message: "Task already exists!",
				Status:  http.StatusConflict,
			}, nil
		}
	}

	// Get task files
	fileNames, err := tsk.ListFiles()
	if err != nil {
		return nil, err
	}

	// Make response
	taskResponse := api_v1.TaskWithFiles{
		Name:      tsk.Name,
		Contract:  tsk.ContractNumber,
		Completed: tsk.Completed,
		Updated:   tsk.Updated,
		Files:     fileNames,
	}
	taskURL := fmt.Sprintf("/api/v1/managers/%s/tasks/%s/%s", manager.Name, tsk.ContractNumber, tsk.Name)
	parsedURL, err := url.Parse(taskURL)
	if err != nil {
		return nil, err
	}
	response := api_v1.TaskWithFilesHeaders{
		Response: taskResponse,
		Location: api_v1.NewOptURI(*parsedURL),
	}

	return &response, nil
}

func (h *TaskHandler) TaskPatch(ctx context.Context, req *api_v1.TaskUpdate, params api_v1.TaskPatchParams) (api_v1.TaskPatchRes, error) {
	manager, ok := h.managers[params.ManagerName]
	if !ok {
		return &api_v1.TaskPatchNotFound{
			Code:    api_v1.ErrorResponseCodeNOTFOUND,
			Message: "Manager not found",
			Status:  http.StatusNotFound,
		}, nil
	}

	// Get task parameters
	currentName := params.Name
	currentContract := params.Contract
	name, hasName := req.Name.Get()
	contract, hasContract := req.Contract.Get()
	completed, hasCompleted := req.Completed.Get()

	// Check parameters
	if !isValidName(currentName) {
		return &api_v1.TaskPatchBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid name!",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if !isValidName(currentContract) {
		return &api_v1.TaskPatchBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid contract!",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if hasContract && !isValidName(contract) {
		return &api_v1.TaskPatchBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid new contract!",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if hasName && !isValidName(name) {
		return &api_v1.TaskPatchBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid new name!",
			Status:  http.StatusBadRequest,
		}, nil
	}

	// Get task
	tsk, err := manager.GetTask(currentContract, currentName)
	if err != nil {
		if errors.Is(err, task.ErrTaskNotExists) {
			return &api_v1.TaskPatchNotFound{
				Code:    api_v1.ErrorResponseCodeNOTFOUND,
				Message: "Task not found!",
				Status:  http.StatusNotFound,
			}, nil
		}
		return nil, err
	}

	switch {
	// Change name and contract
	case hasContract && hasName:
		err := tsk.Move(contract, name)
		if err != nil {
			if errors.Is(err, task.ErrTaskExists) {
				return &api_v1.TaskPatchConflict{
					Code:    api_v1.ErrorResponseCodeCONFLICT,
					Message: "This task already exists!",
					Status:  http.StatusConflict,
				}, nil
			}
			return nil, err
		}
	// Change contract
	case hasContract:
		err := tsk.ChangeContract(contract)
		if err != nil {
			if errors.Is(err, task.ErrTaskExists) {
				return &api_v1.TaskPatchConflict{
					Code:    api_v1.ErrorResponseCodeCONFLICT,
					Message: "Task with that name already exists in this contract!",
					Status:  http.StatusConflict,
				}, nil
			}
			return nil, err
		}
	// Change name
	case hasName:
		err := tsk.Rename(name)
		if err != nil {
			if errors.Is(err, task.ErrTaskExists) {
				return &api_v1.TaskPatchConflict{
					Code:    api_v1.ErrorResponseCodeCONFLICT,
					Message: "Task with that name already exists!",
					Status:  http.StatusConflict,
				}, nil
			}
			return nil, err
		}
	}

	if hasCompleted {
		if completed {
			err = tsk.Complete()
		} else {
			err = tsk.MakeActive()
		}
		if err != nil {
			if errors.Is(err, task.ErrTaskCompleted) {
				return &api_v1.TaskPatchConflict{
					Code:    api_v1.ErrorResponseCodeCONFLICT,
					Message: "Task already completed!",
					Status:  http.StatusConflict,
				}, nil
			}
			if errors.Is(err, task.ErrTaskActive) {
				return &api_v1.TaskPatchConflict{
					Code:    api_v1.ErrorResponseCodeCONFLICT,
					Message: "Task already active!",
					Status:  http.StatusConflict,
				}, nil
			}
			return nil, err
		}
	}

	// Get task files
	files, err := tsk.ListFiles()
	if err != nil {
		return nil, err
	}

	// Make response
	response := api_v1.TaskWithFiles{
		Name:      tsk.Name,
		Contract:  tsk.ContractNumber,
		Completed: tsk.Completed,
		Updated:   tsk.Updated,
		Files:     files,
	}

	return &response, nil
}

func (h *TaskHandler) TaskDelete(ctx context.Context, params api_v1.TaskDeleteParams) (api_v1.TaskDeleteRes, error) {
	manager, ok := h.managers[params.ManagerName]
	if !ok {
		return &api_v1.TaskDeleteNotFound{
			Code:    api_v1.ErrorResponseCodeNOTFOUND,
			Message: "Manager not found",
			Status:  http.StatusNotFound,
		}, nil
	}

	// Get parameters
	contract := params.Contract
	name := params.Name

	//Check parameters
	if !isValidName(contract) {
		return &api_v1.TaskDeleteBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid contract!",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if !isValidName(name) {
		return &api_v1.TaskDeleteBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid name!",
			Status:  http.StatusBadRequest,
		}, nil
	}

	// Get task
	tsk, err := manager.GetTask(contract, name)
	if err != nil {
		if errors.Is(err, task.ErrTaskNotExists) {
			return &api_v1.TaskDeleteNotFound{
				Code:    api_v1.ErrorResponseCodeNOTFOUND,
				Message: "Task not found!",
				Status:  http.StatusNotFound,
			}, nil
		}
		return nil, err
	}

	// Delete task
	err = tsk.Delete()
	if err != nil {
		return nil, err
	}

	return &api_v1.TaskDeleteNoContent{}, nil
}
