package task

import (
	"os"
	"path/filepath"
	"time"
)

func (tsk *Task) Move(contract, name string) error {
	// Check newContract
	if contract == "" {
		return ErrInvalidContract
	}
	// Check newName
	if name == "" {
		return ErrInvalidName
	}

	// Check task with new parameters does not exist
	exists, err := tsk.TaskManager.Exists(contract, name)
	if err != nil {
		return err
	}
	if exists {
		return ErrTaskExists
	}

	// Get old path
	oldPath := tsk.Path()
	// Build new path (stay in same status folder)
	var newPath string
	if tsk.Completed {
		newPath = filepath.Join(tsk.TaskManager.CompletePath, contract, name)
	} else {
		newPath = filepath.Join(tsk.TaskManager.ActivePath, contract, name)
	}

	// Move directory
	err = os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}

	// Update task fields
	tsk.ContractNumber = contract
	tsk.Name = name
	tsk.Updated = time.Now().UTC()
	return nil
}

// ChangeContract изменяет договор задачи.
func (tsk *Task) ChangeContract(newContract string) error {
	return tsk.Move(newContract, tsk.Name)
}

// Rename переименовывает задачу.
func (tsk *Task) Rename(newName string) error {
	return tsk.Move(tsk.ContractNumber, newName)
}
