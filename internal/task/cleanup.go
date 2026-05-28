package task

import (
	"fmt"
	"os"
	"path/filepath"
)

func (m *Manager) CleanupContractIfEmpty(path string) error {
	// List tasks
	taskList, err := m.listTasksContract(path, "", false)
	if err != nil {
		return fmt.Errorf("failed to list tasks: %w", err)
	}

	// Remove if blank
	if len(taskList) == 0 {
		err = os.Remove(path)
		if err != nil {
			return fmt.Errorf("cannot delete void contract: %v", err)
		}
	}

	return nil
}

func (m *Manager) cleanupDir(path string) error {
	// Read dir
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read dir: %w", err)
	}

	// Check subdirectories
	for _, entry := range entries {
		// Check is dir
		if !entry.IsDir() {
			continue
		}

		// Get entry path
		entryPath := filepath.Join(path, entry.Name())

		// Cleanup if empty
		err := m.CleanupContractIfEmpty(entryPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) Cleanup() error {
	err := m.cleanupDir(m.ActivePath)
	if err != nil {
		return fmt.Errorf("failed to cleanup active tasks: %w", err)
	}
	err = m.cleanupDir(m.CompletePath)
	if err != nil {
		return fmt.Errorf("failed to cleanup complete tasks: %w", err)
	}
	return nil
}
