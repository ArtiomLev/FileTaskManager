package task

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileInput struct {
	Name   string
	Reader io.Reader
}

func (m *Manager) CreateTask(contract, name string, files []FileInput) (*Task, error) {
	// Check input
	if contract == "" {
		return nil, ErrInvalidContract
	}
	if name == "" {
		return nil, ErrInvalidName
	}
	if len(files) == 0 {
		return nil, ErrNoFiles
	}

	// Create dir paths
	targetDir := filepath.Join(m.ActivePath, contract, name)
	tempDir := targetDir + ".tmp"

	// Check task exists
	exists, err := m.Exists(contract, name)
	if err != nil {
		return nil, fmt.Errorf("error checking task directory: %w", err)
	}
	if exists {
		return nil, ErrTaskExists
	}

	// Create temp dir
	err = os.MkdirAll(tempDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("error creating temp directory: %w", err)
	}

	// Delete temp dir error
	success := false
	defer func() {
		if !success {
			err := os.RemoveAll(tempDir)
			if err != nil {
				return
			}
		}
	}()

	// Put files to temp dir
	for _, fileInput := range files {
		safeFilename := filepath.Base(fileInput.Name)
		dstPath := filepath.Join(tempDir, safeFilename)
		dst, err := os.Create(dstPath)
		if err != nil {
			return nil, fmt.Errorf("error creating file: %w", err)
		}
		_, err = io.Copy(dst, fileInput.Reader)
		err = dst.Close()
		if err != nil {
			return nil, err
		}
	}

	// Moving temp to task
	if err := os.Rename(tempDir, targetDir); err != nil {
		return nil, ErrCannotMoveTempToTask
	}

	return &Task{
			TaskManager:    m,
			ContractNumber: contract,
			Name:           name,
			Completed:      false,
			Updated:        time.Now().UTC()},
		nil
}
