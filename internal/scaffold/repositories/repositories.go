package repositories

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/lclpedro/scaffold-golang-fiber/internal/scaffold/repositories/health"
	uow "github.com/lclpedro/scaffold-golang-fiber/pkg/unity_of_work"
)

type AllRepositories struct {
	HealthRepository health.HealthRepository
}

func RegistryRepositories(uow uow.UnityOfWorkInterface, dbConnection *sqlx.DB) uow.UnityOfWorkInterface {
	uow.Register("HealthRepository", func(tx *sql.Tx) interface{} {
		repo := health.NewHealthRepository(dbConnection)
		return repo
	})
	return uow
}
