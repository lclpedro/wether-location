package health

import (
	"github.com/jmoiron/sqlx"
)

type HealthRepository interface {
	GetDatabaseCheck() error
}

func NewHealthRepository(dbConnection *sqlx.DB) HealthRepository {
	return &healthRepository{
		dbConnection: dbConnection,
	}
}
