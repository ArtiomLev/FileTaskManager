package api_handler_v1

import (
	"LaserTaskSystem/internal/api/api_v1"
	"context"
	"log"
	"net/http"
)

func (h *ApiHandlerV1) NewError(ctx context.Context, err error) *api_v1.ErrorStatusCode {
	log.Printf("Unexpected error: %v\n", err)
	return &api_v1.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: api_v1.ErrorResponse{
			Code:    api_v1.ErrorResponseCodeINTERNALERROR,
			Message: "Internal server error",
			Status:  http.StatusInternalServerError,
		},
	}
}
