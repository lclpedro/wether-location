package mysql

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

type connection struct {
	read  *sqlx.DB
	write *sqlx.DB
}

var (
	ErrConnectionNotNil = errors.New("write connection could not be nil")
)

type Connection interface {
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	Get(dest interface{}, query string, args ...interface{}) error
	QueryRow(query string, args ...interface{}) *sql.Row
	Close() error
	GetDB() *sqlx.DB
}

func NewConnection(write *sqlx.DB, read *sqlx.DB) (Connection, error) {
	if write == nil {
		return nil, ErrConnectionNotNil
	}
	if read == nil {
		read = write
	}
	return &connection{
		read:  read,
		write: write,
	}, nil
}
func (conn *connection) GetDB() *sqlx.DB {
	return conn.write
}
func (conn *connection) Select(dest interface{}, query string,
	args ...interface{}) (err error) {
	err = conn.read.Select(dest, query, args...)
	return
}

func (conn *connection) Exec(query string,
	args ...interface{}) (result sql.Result, err error) {
	result, err = conn.write.Exec(query, args...)
	return
}

func (conn *connection) NamedExec(query string,
	arg interface{}) (result sql.Result, err error) {
	result, err = conn.write.NamedExec(query, arg)
	return
}
func (conn *connection) NamedQuery(query string,
	arg interface{}) (row *sqlx.Rows, err error) {
	row, err = conn.read.NamedQuery(query, arg)
	return
}

func (conn *connection) Get(dest interface{}, query string,
	args ...interface{}) (err error) {
	err = conn.read.Get(dest, query, args...)
	return
}

func (conn *connection) QueryRow(query string,
	args ...interface{}) (row *sql.Row) {
	row = conn.read.QueryRow(query, args...)
	return
}

func (conn *connection) Close() error {
	err := conn.read.Close()
	if err != nil {
		return err
	}
	err = conn.write.Close()
	return err
}

// mocks
func MockConnection() (sqlmock.Sqlmock, Connection, error) {
	db, mock, _ := sqlmock.New()
	conn := sqlx.NewDb(db, "sqlmock")

	con, err := NewConnection(conn, conn)
	return mock, con, err
}
