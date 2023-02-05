package mysql

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/mercadolibre/fury_treasury-coupons/internal/commons/utils"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type connection struct {
	read  *sqlx.DB
	write *sqlx.DB
}

var (
	ErrConnectionNotNil = errors.New("write connection could not be nil")
)

const product = "DATABASE"

type Connection interface {
	Select(operation string, tnx *newrelic.Transaction, dest interface{}, query string, args ...interface{}) error
	Exec(operation string, txn *newrelic.Transaction, query string, args ...interface{}) (sql.Result, error)
	NamedExec(operation string, txn *newrelic.Transaction, query string, arg interface{}) (sql.Result, error)
	NamedQuery(operation string, txn *newrelic.Transaction, query string, arg interface{}) (*sqlx.Rows, error)
	Get(operation string, txn *newrelic.Transaction, dest interface{}, query string, args ...interface{}) error
	QueryRow(operation string, txn *newrelic.Transaction, query string, args ...interface{}) *sql.Row
	WriteNamedExecTx(operation string, txn *newrelic.Transaction, tx *sqlx.Tx, query string, args ...interface{}) (*sqlx.Tx, sql.Result, error)
	WriteExecTx(operation string, txn *newrelic.Transaction, tx *sqlx.Tx, query string, args ...interface{}) (*sqlx.Tx, sql.Result, error)
	Close() error
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

func (conn *connection) Select(operation string, tnx *newrelic.Transaction, dest interface{}, query string,
	args ...interface{}) (err error) {
	utils.WrapDatastoreSegment(product, operation, tnx, func() {
		err = conn.read.Select(dest, query, args...)
	})
	return
}

func (conn *connection) Exec(operation string, tnx *newrelic.Transaction, query string,
	args ...interface{}) (result sql.Result, err error) {
	utils.WrapDatastoreSegment(product, operation, tnx, func() {
		result, err = conn.write.Exec(query, args...)
	})
	return
}

func (conn *connection) NamedExec(operation string, tnx *newrelic.Transaction, query string,
	arg interface{}) (result sql.Result, err error) {
	utils.WrapDatastoreSegment(product, operation, tnx, func() {
		result, err = conn.write.NamedExec(query, arg)
	})
	return
}
func (conn *connection) NamedQuery(operation string, tnx *newrelic.Transaction, query string,
	arg interface{}) (row *sqlx.Rows, err error) {
	utils.WrapDatastoreSegment(product, operation, tnx, func() {
		row, err = conn.read.NamedQuery(query, arg)
	})
	return
}

func (conn *connection) Get(operation string, tnx *newrelic.Transaction, dest interface{}, query string,
	args ...interface{}) (err error) {
	utils.WrapDatastoreSegment(product, operation, tnx, func() {
		err = conn.read.Get(dest, query, args...)
	})
	return
}

func (conn *connection) QueryRow(operation string, tnx *newrelic.Transaction, query string,
	args ...interface{}) (row *sql.Row) {
	utils.WrapDatastoreSegment(product, operation, tnx, func() {
		row = conn.read.QueryRow(query, args...)
	})
	return
}

// WriteNamedExecTx inits and allows transaction of multiple query executions using the NamedExec process and hold a connection
// until call of `Commit()` or `Rollback()`. More info: https://jmoiron.github.io/sqlx/#transactions
func (conn *connection) WriteNamedExecTx(operation string, tnx *newrelic.Transaction, tx *sqlx.Tx, query string,
	args ...interface{}) (*sqlx.Tx, sql.Result, error) {
	var err error
	if tx == nil {
		tx, err = conn.write.Beginx()
		if err != nil {
			return nil, nil, err
		}
	}

	var result sql.Result
	utils.WrapDatastoreSegment(product, operation, tnx, func() {
		result, err = tx.NamedExec(query, args)
	})
	return tx, result, err
}

// WriteExecTx inits and allows transaction of multiple query executions using the Exec process and hold a connection
// until call of `Commit()` or `Rollback()`. More info: https://jmoiron.github.io/sqlx/#transactions
func (conn *connection) WriteExecTx(operation string, tnx *newrelic.Transaction, tx *sqlx.Tx, query string, args ...interface{}) (*sqlx.Tx, sql.Result, error) {
	var err error
	if tx == nil {
		tx, err = conn.write.Beginx()
		if err != nil {
			return nil, nil, err
		}
	}

	var result sql.Result
	utils.WrapDatastoreSegment(product, operation, tnx, func() {
		result, err = tx.Exec(query, args...)
	})
	return tx, result, err
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
