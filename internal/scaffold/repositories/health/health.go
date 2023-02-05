package health

import (
	"github.com/jmoiron/sqlx"
)

type HealthRepository interface {
	GetDatabaseCheck() error
}

func NewHealthRepository(mysqlConnection *sqlx.DB) HealthRepository {
	return &healthRepository{
		mysqlConnection: mysqlConnection,
	}
}
