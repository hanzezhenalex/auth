//go:build docker

package mysql

import (
	"context"
	"testing"

	"github.com/hanzezhenalex/auth/src"
	"github.com/hanzezhenalex/auth/src/datastore"

	"github.com/stretchr/testify/require"
)

var store *mysqlDatastore

func createMysqlDatastore() *mysqlDatastore {
	const path = "../../../dev/config.json"
	cfg, err := src.NewConfigFromFile(path)
	if err != nil {
		panic(err)
	}

	store, err := NewMysqlDatastore(cfg.DbConfig, true)
	if err != nil {
		panic(err)
	}
	return store
}

func init() {
	//src.EnableDebugMode()
	store = createMysqlDatastore()
}

func TestMysqlDatastore_Authority(t *testing.T) {
	rq := require.New(t)
	ctx := context.Background()

	t.Run("create", func(t *testing.T) {
		const name = "test_1"

		rq.NoError(store.CreateAuthority(ctx, &datastore.Authority{AuthName: name}))
	})

	t.Run("create duplicated authority, should fail", func(t *testing.T) {
		const name = "test_2"

		rq.NoError(store.CreateAuthority(ctx, &datastore.Authority{AuthName: name}))
		rq.Equal(datastore.ErrorAuthExist, store.CreateAuthority(ctx, &datastore.Authority{AuthName: name}))
	})

	t.Run("delete", func(t *testing.T) {
		const name = "test_3"
		auth := &datastore.Authority{AuthName: name}

		rq.NoError(store.CreateAuthority(ctx, auth))
		rq.NoError(store.DeleteAuthorityByID(ctx, auth.ID))
	})

	t.Run("delete a non-existed one, should fail", func(t *testing.T) {
		rq.Equal(datastore.ErrorAuthNotExist, store.DeleteAuthorityByID(ctx, 9999))
	})

	t.Run("create and delete multiple times", func(t *testing.T) {
		const name = "test_4"
		auth := &datastore.Authority{AuthName: name}

		// create and delete first time
		rq.NoError(store.CreateAuthority(ctx, auth))
		rq.NoError(store.DeleteAuthorityByID(ctx, auth.ID))

		// clean up
		auth.ID = 0
		auth.DeletedAt = 0

		// create and delete second time
		rq.NoError(store.CreateAuthority(ctx, auth))
		rq.NoError(store.DeleteAuthorityByID(ctx, auth.ID))
	})
}
