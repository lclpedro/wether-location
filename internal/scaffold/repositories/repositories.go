package repositories

import (
	"database/sql"
	"github.com/lclpedro/scaffold-golang-fiber/internal/scaffold/repositories/health"
	"github.com/lclpedro/scaffold-golang-fiber/pkg/mysql"
	uow "github.com/lclpedro/scaffold-golang-fiber/pkg/unity_of_work"
)

type AllRepositories struct {
	HealthRepository health.HealthRepository
}

func RegistryRepositories(uow uow.UnityOfWorkInterface, dbConnection mysql.Connection) uow.UnityOfWorkInterface {
	uow.Register("HealthRepository", func(tx *sql.Tx) interface{} {
		repo := health.NewHealthRepository(dbConnection)
		return repo
	})
	return uow
}
