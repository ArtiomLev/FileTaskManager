package task

import (
	"fmt"
	"os"
)

func (tsk *Task) Delete() error {
	//Get path
	taskPath := tsk.Path()

	// Remove task
	err := os.RemoveAll(taskPath)
	if err != nil {
		return fmt.Errorf("cannot delete task: %v", err)
	}

	return nil
}
