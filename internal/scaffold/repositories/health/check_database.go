package health

import (
	"github.com/jmoiron/sqlx"
)

type healthRepository struct {
	mysqlConnection *sqlx.DB
}

func (h *healthRepository) GetDatabaseCheck() error {
	_, err := h.mysqlConnection.Exec("SELECT 1;")
	return err
}
