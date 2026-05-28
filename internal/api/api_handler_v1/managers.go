package api_handler_v1

import (
	"LaserTaskSystem/internal/api/api_v1"
	"context"
)

func (h *TaskHandler) ManagersGet(ctx context.Context) (api_v1.ManagersGetRes, error) {
	managers := make(api_v1.Managers, 0, len(h.managers))
	for _, manager := range h.managers {
		managers = append(managers, api_v1.Manager{
			Name:        manager.Name,
			DisplayName: manager.DisplayName,
		})
	}
	return &managers, nil
}
