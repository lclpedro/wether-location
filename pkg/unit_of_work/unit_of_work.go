package uow

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lclpedro/scaffold-golang-fiber/pkg/mysql"
)

type RepositoryFactory func(tx *sqlx.Tx) interface{}

type UnitOfWorkInterface interface {
	Register(name string, repository RepositoryFactory)
	UnRegister(name string)
	GetRepository(ctx context.Context, name string) (interface{}, error)
	Do(ctx context.Context, fn func(uow *UnitOfWork) error) error
	CommitOrRollback() error
}

type UnitOfWork struct {
	ctx          context.Context
	Db           *sqlx.DB
	Tx           *sqlx.Tx
	Repositories map[string]RepositoryFactory
}

const (
	ErrorTxExists     = "UnitOfWork: Transaction already exists"
	ErrorTxNotExists  = "UnitOfWork: No transaction to rollback"
	ErrorExecRollback = "UnitOfWork: Error in execute rollback transaction. Original Error: %s Rollback Error: %s"
	ErrorExecCommit   = "UnitOfWork: Error in execute commit transaction. Original Error: %s Commit Error: %s"
)

func NewUnitOfWork(db mysql.Connection) *UnitOfWork {
	return &UnitOfWork{
		Db:           db.GetDB(),
		Repositories: make(map[string]RepositoryFactory),
	}
}

func (u *UnitOfWork) initTx(ctx context.Context) error {
	if u.Tx != nil {
		return errors.New(ErrorTxExists)
	}
	tx, err := u.Db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	u.Tx = tx
	return nil
}

func (u *UnitOfWork) Register(name string, repository RepositoryFactory) {
	u.Repositories[name] = repository
}
func (u *UnitOfWork) UnRegister(name string) {
	delete(u.Repositories, name)
}
func (u *UnitOfWork) GetRepository(ctx context.Context, name string) (interface{}, error) {
	if u.Tx == nil {
		err := u.initTx(ctx)
		if err != nil {
			return nil, err
		}
	}
	repo := u.Repositories[name](u.Tx)
	return repo, nil
}

func (u *UnitOfWork) Do(ctx context.Context, fn func(uow *UnitOfWork) error) error {
	err := u.initTx(ctx)
	if err != nil {
		return err
	}
	err = fn(u)
	if err != nil {
		if errRollback := u.rollback(); errRollback != nil {
			return errors.New(
				fmt.Sprintf(ErrorExecRollback, err.Error(), errRollback.Error()),
			)
		}
		return err
	}
	return u.CommitOrRollback()
}

func (u *UnitOfWork) CommitOrRollback() error {
	if u.Tx == nil {
		return errors.New(ErrorTxNotExists)
	}
	err := u.Tx.Commit()
	if err != nil {
		if errRollback := u.rollback(); errRollback != nil {
			return errors.New(
				fmt.Sprintf(ErrorExecCommit, err.Error(), errRollback.Error()),
			)
		}
		return err
	}

	u.Tx = nil
	return nil
}

func (u *UnitOfWork) rollback() error {
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
