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
	duplicatedOnPrimaryKey = 1062
)

type mysqlDatastore struct {
	engine *xorm.Engine
}

func NewMysqlDatastore(cfg src.DbConfig) (*mysqlDatastore, error) {
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

	if err := store.onBoarding(); err != nil {
		return nil, fmt.Errorf("fail to onboarding: %w", err)
	}

	return store, nil
}

func (store *mysqlDatastore) tables() []interface{} {
	return []interface{}{
		new(datastore.User),
		new(datastore.Authority),
		new(datastore.Role),
	}
}

func (store *mysqlDatastore) onBoarding() error {
	tables := store.tables()
	if err := store.engine.Sync(tables...); err != nil {
		return fmt.Errorf("fail to sync tables, err=%w", err)
	}
	return nil
}

func (store *mysqlDatastore) cleanup() error {
	tables := store.tables()

	for _, table := range tables {
		name := table.(interface{ TableName() string }).TableName()
		if _, err := store.engine.
			Table(name).
			Where("1=1").
			Unscoped().
			Delete(); err != nil {
			return fmt.Errorf("fail to delete data in table %s, err=%w", name, err)
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
		if mysqlErr.Number == duplicatedOnPrimaryKey {
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

func (store *mysqlDatastore) GetAuthorityByID(ctx context.Context, id int64) (*datastore.Authority, error) {
	auth := datastore.Authority{ID: id}

	if ok, err := store.engine.Context(ctx).Get(&auth); err != nil {
		return nil, err
	} else if !ok {
		return nil, datastore.ErrorAuthNotExist
	}

	return &auth, nil
}
