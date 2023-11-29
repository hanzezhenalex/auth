//go:build docker

package datastore

import (
	"testing"

	"github.com/hanzezhenalex/auth/src"
)

const path = "../../dev/config.json"

func createMysqlDatastore() *mysqlDatastore {
	cfg, err := src.NewConfigFromFile(path)
	if err != nil {
		panic(err)
	}

	store, err := NewMysqlDatastore(cfg.DbConfig)
	if err != nil {
		panic(err)
	}
	return store
}

func TestMysqlDatastore(t *testing.T) {
	createMysqlDatastore()
}
