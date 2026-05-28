package api_handler_v1

import "LaserTaskSystem/internal/task"

type ApiHandlerV1 struct {
	TaskHandler
}

type TaskHandler struct {
	managers map[string]*task.Manager
}

func NewHandler(managers map[string]*task.Manager) *ApiHandlerV1 {
	return &ApiHandlerV1{
		TaskHandler{
			managers: managers,
		},
	}
}
