package task

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func (tsk *Task) Delete() error {
	//Get path
	taskPath := tsk.Path()

	// Remove task
	err := os.RemoveAll(taskPath)
	if err != nil {
		return fmt.Errorf("cannot delete task: %v", err)
	}

	// Cleanup
	contract := filepath.Dir(taskPath)
	err = tsk.TaskManager.CleanupContractIfEmpty(contract)
	if err != nil {
		log.Printf("cannot cleanup contract %s: %v", tsk.ContractNumber, err)
	}

	return nil
}
