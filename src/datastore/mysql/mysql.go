package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/hanzezhenalex/auth/src"
	"github.com/hanzezhenalex/auth/src/datastore"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

const (
	DuplicatedOnPrimaryKey = 1062
)

type mysqlDatastore struct {
	engine *xorm.Engine
}

func NewMysqlDatastore(cfg src.DbConfig, cleanup bool) (*mysqlDatastore, error) {
	engine, err := xorm.NewEngine("mysql", cfg.Dns())
	if err != nil {
		return nil, fmt.Errorf("fail to connect to db: %w", err)
	}

	engine.SetMaxIdleConns(cfg.MaxIdleConns)
	engine.SetMaxOpenConns(cfg.MaxOpenConns)

	if err := engine.Ping(); err != nil {
		return nil, fmt.Errorf("fail to ping db: %w", err)
	}

	store := &mysqlDatastore{engine: engine}

	if err := store.onBoarding(cleanup); err != nil {
		return nil, fmt.Errorf("fail to onboarding: %w", err)
	}

	return store, nil
}

func (store *mysqlDatastore) onBoarding(cleanup bool) error {
	tables := []interface{}{
		new(datastore.User),
		new(datastore.Authority),
		new(datastore.Role),
	}
	if err := store.engine.Sync(tables...); err != nil {
		return fmt.Errorf("fail to sync tables, err=%w", err)
	}
	if cleanup {
		if _, err := store.engine.Delete(tables...); err != nil {
			return fmt.Errorf("fail to delete data in tables, err=%w", err)
		}
	}
	return nil
}

/*
	Authority
*/

func (store *mysqlDatastore) CreateAuthority(ctx context.Context, auth *datastore.Authority) error {
	_, err := store.engine.Context(ctx).Insert(auth)

	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == DuplicatedOnPrimaryKey {
			return datastore.ErrorAuthExist
		}
	}
	return err
}

// DeleteAuthorityByID soft delete
func (store *mysqlDatastore) DeleteAuthorityByID(ctx context.Context, id int64) error {
	auth := &datastore.Authority{
		ID:        id,
		DeletedAt: time.Now().Unix(),
	}
	n, err := store.engine.Context(ctx).
		Cols("deleted_at").
		Where("id=?", auth.ID).
		Update(auth)

	if err != nil {
		return err
	} else if n == 0 {
		return datastore.ErrorAuthNotExist
	} else {
		return nil
	}
}
