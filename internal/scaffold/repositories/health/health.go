package health

import "github.com/jmoiron/sqlx"

type HealthReposotiry interface {
	GetDatabaseCheck() int
}

type healthRepository struct {
	mysqlConnection sqlx.Tx
}
