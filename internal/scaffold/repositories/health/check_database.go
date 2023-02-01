package health

func (h *healthRepository) GetDatabaseCheck() error {
	_, err := h.mysqlConnection.Exec("SELECT 1")
	return err
}
