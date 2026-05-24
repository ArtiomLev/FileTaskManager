package task

import (
	"fmt"
	"os"
	"path/filepath"
)

func (m *Manager) Exists(contract, name string) (bool, error) {
	activePath := filepath.Join(m.ActivePath, contract, name)
	completedPath := filepath.Join(m.CompletePath, contract, name)

	_, activeErr := os.Stat(activePath)
	_, completedErr := os.Stat(completedPath)

	if activeErr == nil || completedErr == nil {
		return true, nil
	}
	if !os.IsNotExist(activeErr) {
		return false, activeErr
	}
	if !os.IsNotExist(completedErr) {
		return false, completedErr
	}
	return false, nil
}

func (m *Manager) listTasks(rootPath string, completed bool) ([]Task, error) {
	var tasks []Task

	contracts, err := os.ReadDir(rootPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	for _, contractEntry := range contracts {
		// Check is not file
		if !contractEntry.IsDir() {
			continue
		}

		// Get contract name
		contractName := contractEntry.Name()

		// Get tasks in contract
		contractPath := filepath.Join(rootPath, contractName)
		contractTasks, err := os.ReadDir(contractPath)
		if err != nil {
			return nil, fmt.Errorf("failed to list tasks: %w", err)
		}

		// Process tasks in contract
		for _, contractTask := range contractTasks {
			// Check is not dir
			if !contractTask.IsDir() {
				continue
			}

			// Get task info
			taskName := contractTask.Name()
			taskDirInfo, err := contractTask.Info()
			if err != nil {
				return nil, fmt.Errorf("failed to get task info: %w", err)
			}
			taskUpdated := taskDirInfo.ModTime().UTC()

			tasks = append(tasks, Task{
				TaskManager:    m,
				ContractNumber: contractName,
				Name:           taskName,
				Completed:      completed,
				Updated:        taskUpdated,
			})
		}
	}

	return tasks, nil
}

func (m *Manager) ListActive() ([]Task, error) {
	return m.listTasks(m.ActivePath, false)
}

func (m *Manager) ListCompleted() ([]Task, error) {
	return m.listTasks(m.CompletePath, true)
}

func (m *Manager) ListAll() ([]Task, error) {
	active, err := m.ListActive()
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	completed, err := m.ListCompleted()
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}
	return append(active, completed...), nil
}

func (m *Manager) GetTask(contract, name string) (*Task, error) {
	activePath := filepath.Join(m.ActivePath, contract, name)
	if info, err := os.Stat(activePath); err == nil {
		return &Task{
			TaskManager:    m,
			ContractNumber: contract,
			Name:           name,
			Completed:      false,
			Updated:        info.ModTime().UTC(),
		}, nil
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	completePath := filepath.Join(m.CompletePath, contract, name)
	if info, err := os.Stat(completePath); err == nil {
		return &Task{
			TaskManager:    m,
			ContractNumber: contract,
			Name:           name,
			Completed:      true,
			Updated:        info.ModTime().UTC(),
		}, nil
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return nil, ErrTaskNotExists
}
