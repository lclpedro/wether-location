package mysql

import (
	"github.com/spf13/viper"
	"os"
	"time"
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
	for _, op := range operations {
		databasesConfig[op] = &DatabaseConfig{
			User:               viper.GetString("database." + op + ".username"),
			Host:               os.Getenv(viper.GetString("database." + op + ".hostname")),
			Password:           os.Getenv(viper.GetString("database." + op + ".password")),
			DbName:             viper.GetString("database." + op + ".name"),
			MaxOpenConnections: viper.GetInt("database." + op + ".max_connection_open"),
			MaxIdleConnections: viper.GetInt("database." + op + ".max_connection_idle"),
			MaxConnLifetime:    time.Duration(viper.GetInt("database."+op+".max_connection_life_time")) * time.Minute,
		}
	}

	return databasesConfig
}
