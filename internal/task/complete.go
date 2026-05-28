package task

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func (tsk *Task) moveTask(src, dst string, completed bool) error {
	err := os.MkdirAll(filepath.Dir(dst), 0775)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %v", err)
	}

	err = os.Rename(src, dst)
	if err != nil {
		return fmt.Errorf("failed to move task: %v", err)
	}

	tsk.Updated = time.Now().UTC()
	tsk.Completed = completed
	return nil
}

func (tsk *Task) Complete() error {
	if tsk.Completed == true {
		return ErrTaskCompleted
	}

	src := filepath.Join(tsk.TaskManager.ActivePath, tsk.ContractNumber, tsk.Name)
	dst := filepath.Join(tsk.TaskManager.CompletePath, tsk.ContractNumber, tsk.Name)

	return tsk.moveTask(src, dst, true)
}

func (tsk *Task) MakeActive() error {
	if tsk.Completed == false {
		return ErrTaskActive
	}

	src := filepath.Join(tsk.TaskManager.CompletePath, tsk.ContractNumber, tsk.Name)
	dst := filepath.Join(tsk.TaskManager.ActivePath, tsk.ContractNumber, tsk.Name)

	return tsk.moveTask(src, dst, false)
}
