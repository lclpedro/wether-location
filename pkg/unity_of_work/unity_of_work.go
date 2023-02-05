package uow

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"

	"errors"
)

type RepositoryFactory func(tx *sql.Tx) interface{}

type UnityOfWorkInterface interface {
	Register(name string, repository RepositoryFactory)
	UnRegister(name string)
	GetRepository(ctx context.Context, name string) (interface{}, error)
	Do(ctx context.Context, fn func(uow *UnityOfWork) error) error
	CommitOrRollback() error
}

type UnityOfWork struct {
	ctx          context.Context
	Db           *sqlx.DB
	Tx           *sql.Tx
	Repositories map[string]RepositoryFactory
}

const (
	ErrorTxExists     = "UnityOfWork: Transaction already exists"
	ErrorTxNotExists  = "UnityOfWork: No transaction to rollback"
	ErrorExecRollback = "UnityOfWork: Error in execute rollback transaction. Original Error: %s Rollback Error: %s"
	ErrorExecCommit   = "UnityOfWork: Error in execute commit transaction. Original Error: %s Commit Error: %s"
)

func NewUnityOfWork(db *sqlx.DB) *UnityOfWork {
	return &UnityOfWork{
		Db:           db,
		Repositories: make(map[string]RepositoryFactory),
	}
}

func (u *UnityOfWork) initTx(ctx context.Context) error {
	if u.Tx != nil {
		return errors.New(ErrorTxExists)
	}
	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	u.Tx = tx
	return nil
}

func (u *UnityOfWork) Register(name string, repository RepositoryFactory) {
	u.Repositories[name] = repository
}
func (u *UnityOfWork) UnRegister(name string) {
	delete(u.Repositories, name)
}
func (u *UnityOfWork) GetRepository(ctx context.Context, name string) (interface{}, error) {
	if u.Tx == nil {
		err := u.initTx(ctx)
		if err != nil {
			return nil, err
		}
	}
	repo := u.Repositories[name](u.Tx)
	return repo, nil
}

func (u *UnityOfWork) Do(ctx context.Context, fn func(uow *UnityOfWork) error) error {
	err := u.initTx(ctx)
	if err != nil {
		return err
	}
	err = fn(u)
	if err != nil {
		if errRoolback := u.rollback(); errRoolback != nil {
			return errors.New(
				fmt.Sprintf(ErrorExecRollback, err.Error(), errRoolback.Error()),
			)
		}
		return err
	}
	return u.CommitOrRollback()
}

func (u *UnityOfWork) CommitOrRollback() error {
	if u.Tx == nil {
		return errors.New(ErrorTxNotExists)
	}
	err := u.Tx.Commit()
	if err != nil {
		if errRoolback := u.rollback(); errRoolback != nil {
			return errors.New(
				fmt.Sprintf(ErrorExecCommit, err.Error(), errRoolback.Error()),
			)
		}
		return err
	}

	u.Tx = nil
	return nil
}

func (u *UnityOfWork) rollback() error {
	if u.Tx == nil {
		return errors.New(ErrorTxNotExists)
	}
	err := u.Tx.Rollback()
	if err != nil {
		return err
	}
	u.Tx = nil
	return nil
}
