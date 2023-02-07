package health

func (h *healthRepository) GetDatabaseCheck() error {
	_, err := h.dbConnection.Exec("SELECT 1;")
	return err
}
