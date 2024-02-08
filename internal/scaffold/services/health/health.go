package health

import (
	"context"
)

type Service interface {
	Ping(ctx context.Context) string
}

type healthService struct {
}

func NewHealthService() Service {
	return &healthService{}
}

func (h *healthService) Ping(ctx context.Context) string {
	return "Ok"
}
