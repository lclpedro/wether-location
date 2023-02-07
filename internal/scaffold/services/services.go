package services

import (
	"github.com/lclpedro/scaffold-golang-fiber/internal/scaffold/services/health"
	uow "github.com/lclpedro/scaffold-golang-fiber/pkg/unity_of_work"
)

type AllServices struct {
	HealthService health.Service
}

func NewAllServices(uow uow.UnityOfWorkInterface) *AllServices {
	return &AllServices{
		HealthService: health.NewHealthService(uow),
	}
}
