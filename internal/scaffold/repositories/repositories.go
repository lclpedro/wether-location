package repositories

import (
	"github.com/jmoiron/sqlx"
	"github.com/lclpedro/scaffold-golang-fiber/internal/scaffold/repositories/health"
	"github.com/lclpedro/scaffold-golang-fiber/pkg/mysql"
	uow "github.com/lclpedro/scaffold-golang-fiber/pkg/unit_of_work"
)

type AllRepositories struct {
	HealthRepository health.Repository
}

func RegistryRepositories(uow uow.UnitOfWorkInterface, dbConnection mysql.Connection) uow.UnitOfWorkInterface {
	uow.Register("HealthRepository", func(tx *sqlx.Tx) interface{} {
		repo := health.NewHealthRepository(dbConnection)
		return repo
	})
	return uow
}
