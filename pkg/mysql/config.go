package mysql

import (
	"time"

	"github.com/magiconair/properties"
	"github.com/mercadolibre/fury_go-platform/pkg/fury/secret"
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

func GetDatabaseConfiguration(cfg *properties.Properties) map[string]*DatabaseConfig {
	databasesConfig := make(map[string]*DatabaseConfig)

	operations := [2]string{ReadOperation, WriteOperation}
	for _, op := range operations {
		databasesConfig[op] = &DatabaseConfig{
			User:               cfg.GetString("database."+op+".username", "root"),
			Host:               secret.FromEnv(cfg.GetString("database."+op+".hostname.key", "DB_MYSQL_LOCAL_HOSTNAME")),
			Password:           secret.FromEnv(cfg.GetString("database."+op+".password.key", "DB_MYSQL_LOCAL_PWD")),
			DbName:             cfg.GetString("database."+op+".name", "couponsp"),
			MaxOpenConnections: cfg.GetInt("database."+op+".max_open_conn", 1),
			MaxIdleConnections: cfg.GetInt("database."+op+".max_idle_conn", 1),
			MaxConnLifetime:    time.Duration(cfg.GetInt("database."+op+".max_conn_lifetime.minutes", 30)) * time.Minute,
		}
	}

	return databasesConfig
}
