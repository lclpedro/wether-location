package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func InitMySQLConnection(dbConfig *DatabaseConfig, operation string) *sqlx.DB {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.DbName,
	)
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
	conn.Exec("select 1;")
	// add this to the factory reference
	return conn
}
