package task

import "path/filepath"

func (tsk *Task) Path() string {
	if tsk.Completed {
		return filepath.Join(tsk.TaskManager.CompletePath, tsk.ContractNumber, tsk.Name)
	}

	return filepath.Join(tsk.TaskManager.ActivePath, tsk.ContractNumber, tsk.Name)
}
