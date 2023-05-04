package mysql

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	ReadOperation  = "read"
	WriteOperation = "write"
)

type DatabaseConfig struct {
	User               string
	Password           string
	Host               string
	DbName             string
	MaxOpenConnections int
	MaxIdleConnections int
	MaxConnLifetime    time.Duration
}

func GetDatabaseConfiguration() map[string]*DatabaseConfig {
	databasesConfig := make(map[string]*DatabaseConfig)
	operations := [2]string{ReadOperation, WriteOperation}
	for _, operation := range operations {
		databasesConfig[operation] = &DatabaseConfig{
			User:               viper.GetString("database." + operation + ".username"),
			Host:               os.Getenv(viper.GetString("database." + operation + ".hostname")),
			Password:           os.Getenv(viper.GetString("database." + operation + ".password")),
			DbName:             viper.GetString("database." + operation + ".name"),
			MaxOpenConnections: viper.GetInt("database." + operation + ".max_connection_open"),
			MaxIdleConnections: viper.GetInt("database." + operation + ".max_connection_idle"),
			MaxConnLifetime:    time.Duration(viper.GetInt("database."+operation+".max_connection_life_time")) * time.Minute,
		}
	}

	return databasesConfig
}
