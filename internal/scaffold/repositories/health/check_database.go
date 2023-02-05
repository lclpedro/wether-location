package health

import (
	"github.com/lclpedro/scaffold-golang-fiber/pkg/mysql"
)

type healthRepository struct {
	dbConnection mysql.Connection
}

func (h *healthRepository) GetDatabaseCheck() error {
	_, err := h.dbConnection.Exec("SELECT 1;")
	return err
}
