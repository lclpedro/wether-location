package health

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type healthRepository struct {
	mysqlConnection *sqlx.DB
}

func (h *healthRepository) GetDatabaseCheck() error {
	result, err := h.mysqlConnection.Exec("SELECT now();")
	fmt.Println(result.RowsAffected())
	return err
}
