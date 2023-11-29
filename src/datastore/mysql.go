package datastore

import (
	"fmt"

	"github.com/hanzezhenalex/auth/src"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
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

func (store *mysqlDatastore) onBoarding() error {
	return store.engine.Sync(
		new(User),
		new(Authority),
		new(Role),
	)
}
