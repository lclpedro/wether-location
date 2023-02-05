package health

import (
	"github.com/jmoiron/sqlx"
)

type healthRepository struct {
	dbConnection *sqlx.DB
}

func (h *healthRepository) GetDatabaseCheck() error {
	_, err := h.dbConnection.Exec("SELECT 1;")
	return err
}
