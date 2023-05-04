package repositories

import (
	"database/sql"

	"github.com/lclpedro/scaffold-golang-fiber/internal/scaffold/repositories/health"
	"github.com/lclpedro/scaffold-golang-fiber/pkg/mysql"
)

type AllRepositories struct {
	HealthRepository health.Repository
}

func RegistryRepositories(uow mysql.UnitOfWorkInterface, dbConnection mysql.Connection) mysql.UnitOfWorkInterface {
	uow.Register("HealthRepository", func(tx *sql.Tx) interface{} {
		repo := health.NewHealthRepository(dbConnection)
		return repo
	})
	return uow
}
