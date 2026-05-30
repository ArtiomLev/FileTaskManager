package api_handler_v1

import (
	"LaserTaskSystem/internal/api/api_v1"
	"LaserTaskSystem/internal/task"
	"context"
	"errors"
	"net/http"
)

func (h *TaskHandler) TaskFileGet(ctx context.Context, params api_v1.TaskFileGetParams) (api_v1.TaskFileGetRes, error) {
	manager, ok := h.managers[params.ManagerName]
	if !ok {
		return &api_v1.TaskFileGetNotFound{
			Code:    api_v1.ErrorResponseCodeNOTFOUND,
			Message: "Manager not found",
			Status:  http.StatusNotFound,
		}, nil
	}

	// Get parameters
	contract := params.Contract
	name := params.Name
	filename := params.Filename

	//Check parameters
	if !isValidName(contract) {
		return &api_v1.TaskFileGetBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid contract!",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if !isValidName(name) {
		return &api_v1.TaskFileGetBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid name!",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if !isValidName(filename) {
		return &api_v1.TaskFileGetBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid filename!",
			Status:  http.StatusBadRequest,
		}, nil
	}

	tsk, err := manager.GetTask(contract, name)
	if err != nil {
		if errors.Is(err, task.ErrTaskNotExists) {
			return &api_v1.TaskFileGetNotFound{
				Code:    api_v1.ErrorResponseCodeNOTFOUND,
				Message: "Task not exists!",
				Status:  http.StatusNotFound,
			}, nil
		}
		return nil, err
	}

	file, err := tsk.ReadFile(filename)
	if err != nil {
		if errors.Is(err, task.ErrFileNotExists) {
			return &api_v1.TaskFileGetNotFound{
				Code:    api_v1.ErrorResponseCodeNOTFOUND,
				Message: "File not found in task!",
				Status:  http.StatusNotFound,
			}, nil
		}
		return nil, err
	}

	response := api_v1.TaskFileGetOK{
		Data: file,
	}

	return &response, nil
}

func (h *TaskHandler) TaskFilePost(ctx context.Context, req api_v1.TaskFilePostReq, params api_v1.TaskFilePostParams) (api_v1.TaskFilePostRes, error) {
	manager, ok := h.managers[params.ManagerName]
	if !ok {
		return &api_v1.TaskFilePostNotFound{
			Code:    api_v1.ErrorResponseCodeNOTFOUND,
			Message: "Manager not found",
			Status:  http.StatusNotFound,
		}, nil
	}

	// Get parameters
	contract := params.Contract
	name := params.Name
	filename := params.Filename

	//Check parameters
	if !isValidName(contract) {
		return &api_v1.TaskFilePostBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid contract!",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if !isValidName(name) {
		return &api_v1.TaskFilePostBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid name!",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if !isValidName(filename) {
		return &api_v1.TaskFilePostBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid filename!",
			Status:  http.StatusBadRequest,
		}, nil
	}

	tsk, err := manager.GetTask(contract, name)
	if err != nil {
		if errors.Is(err, task.ErrTaskNotExists) {
			return &api_v1.TaskFilePostNotFound{
				Code:    api_v1.ErrorResponseCodeNOTFOUND,
				Message: "Task not exists!",
				Status:  http.StatusNotFound,
			}, nil
		}
		return nil, err
	}

	err = tsk.AddFile(filename, req.Data, false)
	if err != nil {
		if errors.Is(err, task.ErrFileExists) {
			return &api_v1.TaskFilePostConflict{
				Code:    api_v1.ErrorResponseCodeCONFLICT,
				Message: "File already exists!",
				Status:  http.StatusConflict,
			}, nil
		}
		return nil, err
	}

	fileList, err := tsk.ListFiles()
	if err != nil {
		return nil, err
	}

	return new(api_v1.Files(fileList)), nil
}

func (h *TaskHandler) TaskFilePut(ctx context.Context, req api_v1.TaskFilePutReq, params api_v1.TaskFilePutParams) (api_v1.TaskFilePutRes, error) {
	manager, ok := h.managers[params.ManagerName]
	if !ok {
		return &api_v1.TaskFilePutNotFound{
			Code:    api_v1.ErrorResponseCodeNOTFOUND,
			Message: "Manager not found",
			Status:  http.StatusNotFound,
		}, nil
	}

	// Get parameters
	contract := params.Contract
	name := params.Name
	filename := params.Filename

	//Check parameters
	if !isValidName(contract) {
		return &api_v1.TaskFilePutBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid contract!",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if !isValidName(name) {
		return &api_v1.TaskFilePutBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid name!",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if !isValidName(filename) {
		return &api_v1.TaskFilePutBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid filename!",
			Status:  http.StatusBadRequest,
		}, nil
	}

	tsk, err := manager.GetTask(contract, name)
	if err != nil {
		if errors.Is(err, task.ErrTaskNotExists) {
			return &api_v1.TaskFilePutNotFound{
				Code:    api_v1.ErrorResponseCodeNOTFOUND,
				Message: "Task not exists!",
				Status:  http.StatusNotFound,
			}, nil
		}
		return nil, err
	}

	err = tsk.AddFile(filename, req.Data, true)
	if err != nil {
		return nil, err
	}

	fileList, err := tsk.ListFiles()
	if err != nil {
		return nil, err
	}

	return new(api_v1.Files(fileList)), nil
}

func (h *TaskHandler) TaskFileDelete(ctx context.Context, params api_v1.TaskFileDeleteParams) (api_v1.TaskFileDeleteRes, error) {
	manager, ok := h.managers[params.ManagerName]
	if !ok {
		return &api_v1.TaskFileDeleteNotFound{
			Code:    api_v1.ErrorResponseCodeNOTFOUND,
			Message: "Manager not found",
			Status:  http.StatusNotFound,
		}, nil
	}

	// Get parameters
	contract := params.Contract
	name := params.Name
	filename := params.Filename

	//Check parameters
	if !isValidName(contract) {
		return &api_v1.TaskFileDeleteBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid contract!",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if !isValidName(name) {
		return &api_v1.TaskFileDeleteBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid name!",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if !isValidName(filename) {
		return &api_v1.TaskFileDeleteBadRequest{
			Code:    api_v1.ErrorResponseCodeBADREQUEST,
			Message: "Invalid filename!",
			Status:  http.StatusBadRequest,
		}, nil
	}

	tsk, err := manager.GetTask(contract, name)
	if err != nil {
		if errors.Is(err, task.ErrTaskNotExists) {
			return &api_v1.TaskFileDeleteNotFound{
				Code:    api_v1.ErrorResponseCodeNOTFOUND,
				Message: "Task not exists!",
				Status:  http.StatusNotFound,
			}, nil
		}
		return nil, err
	}

	err = tsk.RemoveFile(filename)
	if err != nil {
		if errors.Is(err, task.ErrFileNotExists) {
			return &api_v1.TaskFileDeleteNotFound{
				Code:    api_v1.ErrorResponseCodeNOTFOUND,
				Message: "File not found in task!",
				Status:  http.StatusNotFound,
			}, nil
		}
		return nil, err
	}

	return &api_v1.TaskFileDeleteNoContent{}, nil
}
