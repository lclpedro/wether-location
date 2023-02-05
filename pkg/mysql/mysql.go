package mysql

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/mercadolibre/fury_go-core/pkg/log"
	"github.com/mercadolibre/fury_treasury-coupons/internal/commons/utils"
)

func InitMySQLConnection(dbConfig *DatabaseConfig, operation string) *sqlx.DB {

	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.DbName)
	conn, err := sqlx.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	conn.DB.SetMaxOpenConns(dbConfig.MaxOpenConnections)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	conn.DB.SetMaxIdleConns(dbConfig.MaxIdleConnections)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	conn.DB.SetConnMaxLifetime(dbConfig.MaxConnLifetime)

	healthFunction()(conn, dbConfig.Host, dbConfig.DbName, operation)

	// add this to the factory reference
	return conn
}

func healthFunction() func(*sqlx.DB, string, string, string) {
	const databaseCheckSleepDuration = 2 * time.Second
	const databaseCheckMaxRetries = 5

	return func(c *sqlx.DB, host, dbName, operation string) {
		ctx := context.Background()
		err := utils.Retry(c.Ping, databaseCheckMaxRetries, databaseCheckSleepDuration)
		if err != nil {
			log.Panic(ctx, fmt.Sprintf("Not able to establish connection to database. Error: %s Host: %s DBName: %s", err.Error(), host, dbName))
		}

		for checkAttempts := 1; checkAttempts <= databaseCheckMaxRetries; checkAttempts++ {
			_, err = c.Query("SELECT 1")
			if err != nil {
				log.Error(ctx, fmt.Sprintf(
					"[retry:%d] Not able to use connection and query database. Error: %s Host: %s DBName: %s MaxRetryAttempts: %d",
					checkAttempts,
					err.Error(),
					host,
					dbName,
					databaseCheckMaxRetries,
				))

				if checkAttempts == databaseCheckMaxRetries {
					log.Panic(ctx, fmt.Sprintf(
						"[PANIC] Not able to establish connection to database. Error: %s Host: %s DBName: %s",
						err.Error(),
						host,
						dbName,
					))
				}
				time.Sleep(databaseCheckSleepDuration)
			} else {
				break
			}
		}

		log.Info(ctx, fmt.Sprintf("Database connection stats: [type:%s][%+v]", operation, c.Stats()))
	}
}
