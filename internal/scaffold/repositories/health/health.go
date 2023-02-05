package health

import (
	"github.com/lclpedro/scaffold-golang-fiber/pkg/mysql"
)

type HealthRepository interface {
	GetDatabaseCheck() error
}

func NewHealthRepository(dbConnection mysql.Connection) HealthRepository {
	return &healthRepository{
		dbConnection: dbConnection,
	}
}
