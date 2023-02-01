package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/lclpedro/scaffold-golang-fiber/internal/scaffold/repositories/health"
)

type AllRepositories struct {
	HealthReposotiry health.HealthReposotiry
}

func NewAllRepositories(conn sqlx.Tx) *AllRepositories {
	return &AllRepositories{}
}
