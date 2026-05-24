package task

import (
	"fmt"
	"os"
)

func (tsk *Task) Delete() error {
	taskPath := tsk.Path()
	err := os.RemoveAll(taskPath)
	if err != nil {
		return fmt.Errorf("cannot delete task: %v", err)
	}
	return nil
}
