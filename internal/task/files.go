package task

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

var forbiddenNames = map[string]bool{
	"":          true,
	".":         true,
	"..":        true,
	"meta.yaml": true,
}

func checkFilename(name string) error {
	if forbiddenNames[name] {
		return ErrFileName
	}
	return nil
}

func (tsk *Task) ListFiles() ([]string, error) {
	taskDir, err := os.ReadDir(tsk.Path())
	if err != nil {
		return []string{}, fmt.Errorf("error listing files: %w", err)
	}

	var files []string
	for _, entry := range taskDir {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		if fileName == "meta.yaml" {
			continue
		}

		files = append(files, fileName)
	}

	return files, nil
}

func (tsk *Task) AddFile(name string, reader io.Reader, replace bool) error {
	fileName := filepath.Base(name)
	finalPath := filepath.Join(tsk.Path(), fileName)

	err := checkFilename(fileName)
	if err != nil {
		return err
	}

	if _, err := os.Stat(finalPath); err == nil {
		if !replace {
			return ErrFileExists
		}
		if err := os.Remove(finalPath); err != nil {
			return fmt.Errorf("cannot remove existing file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("cannot stat file: %w", err)
	}

	tempName := fmt.Sprintf(".tmp_%s_%d", fileName, time.Now().UnixNano())
	tempPath := filepath.Join(tsk.Path(), tempName)

	tmpFile, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("cannot create temp file: %w", err)
	}

	success := false
	defer func() {
		_ = tmpFile.Close()
		if !success {
			_ = os.Remove(tempPath)
		}
	}()

	if _, err := io.Copy(tmpFile, reader); err != nil {
		return fmt.Errorf("copy failed: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("close temp file failed: %w", err)
	}

	if err := os.Rename(tempPath, finalPath); err != nil {
		return fmt.Errorf("rename failed: %w", err)
	}

	success = true
	tsk.Updated = time.Now().UTC()
	return nil
}

func (tsk *Task) RemoveFile(name string) error {
	fileName := filepath.Base(name)
	filePath := filepath.Join(tsk.Path(), fileName)

	err := checkFilename(fileName)
	if err != nil {
		return err
	}

	err = os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrFileNotExists
		}
		return err
	}

	tsk.Updated = time.Now().UTC()
	return nil
}

func (tsk *Task) ReadFile(name string) (*os.File, error) {
	fileName := filepath.Base(name)
	filePath := filepath.Join(tsk.Path(), fileName)

	err := checkFilename(fileName)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}

	return file, nil
}
