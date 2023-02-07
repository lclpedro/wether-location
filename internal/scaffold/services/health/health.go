package health

import (
	"context"
	"github.com/lclpedro/scaffold-golang-fiber/internal/scaffold/repositories/health"
	uow "github.com/lclpedro/scaffold-golang-fiber/pkg/unity_of_work"
)

type Service interface {
	Ping(ctx context.Context) error
}

type healthService struct {
	uow uow.UnityOfWorkInterface
}

func NewHealthService(uow uow.UnityOfWorkInterface) Service {
	return &healthService{
		uow: uow,
	}
}

func (h *healthService) getHealthRepository(ctx context.Context) (health.Repository, error) {
	repo, err := h.uow.GetRepository(ctx, "HealthRepository")
	if err != nil {
		return nil, err
	}
	return repo.(health.Repository), nil
}

func (h *healthService) Ping(ctx context.Context) error {
	return h.uow.Do(ctx, func(uow *uow.UnityOfWork) error {
		repo, err := h.getHealthRepository(ctx)

		if err != nil {
			return err
		}

		err = repo.GetDatabaseCheck()

		if err != nil {
			return err
		}

		return nil
	})
}
