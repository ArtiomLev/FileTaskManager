package task

import "time"

type Manager struct {
	Name         string `json:"name"`
	DisplayName  string `json:"display_name"`
	ActivePath   string `json:"-"`
	CompletePath string `json:"-"`
}

func NewManager(name, displayName, activePath, completePath string) (*Manager, error) {
	return &Manager{
			Name:         name,
			DisplayName:  displayName,
			ActivePath:   activePath,
			CompletePath: completePath},
		nil
}

type Task struct {
	TaskManager    *Manager  `json:"-"`
	ContractNumber string    `json:"contract_number"`
	Name           string    `json:"name"`
	Completed      bool      `json:"completed"`
	Updated        time.Time `json:"updated"`
}
