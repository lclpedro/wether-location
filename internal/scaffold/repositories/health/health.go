package health

import (
	"github.com/lclpedro/scaffold-golang-fiber/pkg/mysql"
)

type Repository interface {
	GetDatabaseCheck() error
}

type healthRepository struct {
	dbConnection mysql.Connection
}

func NewHealthRepository(dbConnection mysql.Connection) Repository {
	return &healthRepository{
		dbConnection: dbConnection,
	}
}
