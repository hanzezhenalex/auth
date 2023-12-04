//go:build docker

package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/hanzezhenalex/auth/src"
	"github.com/hanzezhenalex/auth/src/datastore"

	"github.com/stretchr/testify/require"
)

var store *mysqlDatastore

const nonExistedID = 99999

func createMysqlDatastore() *mysqlDatastore {
	const path = "../../../dev/config.json"
	cfg, err := src.NewConfigFromFile(path)
	if err != nil {
		panic(err)
	}

	store, err := NewMysqlDatastore(cfg.DbConfig)
	if err != nil {
		panic(err)
	}
	if err := store.cleanup(); err != nil {
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
		rq.NoError(store.DeleteAuthorityByID(ctx, auth.ID, false))
	})

	t.Run("delete a non-existed one, should fail", func(t *testing.T) {
		rq.Equal(datastore.ErrorAuthNotExist, store.DeleteAuthorityByID(ctx, 9999, false))
	})

	t.Run("read", func(t *testing.T) {
		const name = "test_5"
		expected := &datastore.Authority{AuthName: name}

		rq.NoError(store.CreateAuthority(ctx, expected))

		actual, err := store.GetAuthorityByID(ctx, expected.ID)

		rq.NoError(err)
		rq.Equal(expected.ID, actual.ID)
		rq.Equal(expected.AuthName, actual.AuthName)
		rq.Equal(expected.CreatedAt.Unix(), actual.CreatedAt.Unix())
	})

	t.Run("read a non-existed one, should fail", func(t *testing.T) {
		actual, err := store.GetAuthorityByID(ctx, 9999)

		rq.Nil(actual)
		rq.Equal(datastore.ErrorAuthNotExist, err)
	})

	t.Run("read a deleted one", func(t *testing.T) {
		const name = "test_auth_read_a_deleted_one"
		expected := &datastore.Authority{AuthName: name}

		rq.NoError(store.CreateAuthority(ctx, expected))
		rq.NoError(store.DeleteAuthorityByID(ctx, expected.ID, false))

		actual, err := store.GetAuthorityByID(ctx, expected.ID)

		rq.Nil(actual)
		rq.Equal(datastore.ErrorAuthNotExist, err)
	})

	t.Run("create and delete multiple times", func(t *testing.T) {
		const name = "test_4"
		auth := &datastore.Authority{AuthName: name}

		// create and delete first time
		rq.NoError(store.CreateAuthority(ctx, auth))
		rq.NoError(store.DeleteAuthorityByID(ctx, auth.ID, false))

		// clean up
		auth.ID = 0
		auth.DeletedAt = 0

		// create and delete second time
		rq.NoError(store.CreateAuthority(ctx, auth))
		rq.NoError(store.DeleteAuthorityByID(ctx, auth.ID, false))
	})

	t.Run("delete an auth with role binding", func(t *testing.T) {
		const authName = "test_auth_delete_an_auth_with_role_binding"
		expected := &datastore.Authority{AuthName: authName}

		rq.NoError(store.CreateAuthority(ctx, expected))

		role := datastore.Role{
			RoleName: "test_auth_delete_an_auth_with_role_binding_role",
			Scopes:   []string{"scope1", "scope2"},
			Auths:    []string{authName},
		}
		rq.NoError(store.CreateRole(ctx, &role))

		rq.Equal(
			datastore.ErrorDeleteAuthWithBinding,
			store.DeleteAuthorityByID(ctx, expected.ID, false),
		)

		rq.NoError(store.DeleteAuthorityByID(ctx, expected.ID, true))
	})
}

func TestMysqlDatastore_Role(t *testing.T) {
	rq := require.New(t)
	ctx := context.Background()

	const auth1 = "test_role_auth_1"
	rq.NoError(store.CreateAuthority(ctx, &datastore.Authority{AuthName: auth1}))

	const auth2 = "test_role_auth_2"
	rq.NoError(store.CreateAuthority(ctx, &datastore.Authority{AuthName: auth2}))

	const auth3 = "test_role_auth_3_not_exist"

	t.Run("create, no auth", func(t *testing.T) {
		role := datastore.Role{
			RoleName: "test_role_create_no_auth",
			Scopes:   []string{"scope1", "scope2"},
		}
		rq.NoError(store.CreateRole(ctx, &role))
	})

	t.Run("create, with auth", func(t *testing.T) {
		role := datastore.Role{
			RoleName: "test_role_create_with_auth",
			Scopes:   []string{"scope1", "scope2"},
			Auths:    []string{auth1, auth2},
		}
		rq.NoError(store.CreateRole(ctx, &role))
	})

	t.Run("create, with non-existed auth", func(t *testing.T) {
		role := datastore.Role{
			RoleName: "test_role_create_with non-existed auth",
			Scopes:   []string{"scope1", "scope2"},
			Auths:    []string{auth1, auth3},
		}
		rq.Equal(datastore.ErrorAuthNotExist, store.CreateRole(ctx, &role))
	})

	t.Run("delete", func(t *testing.T) {
		role := datastore.Role{
			RoleName: "test_role_delete",
			Scopes:   []string{"scope1", "scope2"},
			Auths:    []string{auth1, auth2},
		}
		rq.NoError(store.CreateRole(ctx, &role))

		rq.NoError(store.DeleteRoleByID(ctx, role.ID))
	})

	t.Run("delete an non-existed role", func(t *testing.T) {
		rq.Equal(datastore.ErrorRoleNotExist, store.DeleteRoleByID(ctx, nonExistedID))
	})

	t.Run("read", func(t *testing.T) {
		expected := &datastore.Role{
			RoleName: "test_role_read",
			Scopes:   []string{"scope1", "scope2"},
			Auths:    []string{auth1, auth2},
		}
		rq.NoError(store.CreateRole(ctx, expected))

		actual, err := store.GetRoleByID(ctx, expected.ID)
		rq.NoError(err)

		expected.CreatedAt = time.Unix(expected.CreatedAt.Unix(), 0)
		actual.CreatedAt = time.Unix(actual.CreatedAt.Unix(), 0)

		rq.EqualValues(expected, actual)
	})

	t.Run("read an non-existed role", func(t *testing.T) {
		_, err := store.GetRoleByID(ctx, nonExistedID)
		rq.Equal(datastore.ErrorRoleNotExist, err)
	})

	t.Run("read a role with deleted auth", func(t *testing.T) {
		const deletedAuth = "read_a_role_with_deleted_auth"
		auth := &datastore.Authority{AuthName: deletedAuth}
		rq.NoError(store.CreateAuthority(ctx, auth))

		role := &datastore.Role{
			RoleName: "test_role_read_a_role_with_deleted_auth",
			Scopes:   []string{"scope1", "scope2"},
			Auths:    []string{auth1, deletedAuth},
		}

		rq.NoError(store.CreateRole(ctx, role))

		rq.NoError(store.DeleteAuthorityByID(ctx, auth.ID, true))

		actual, err := store.GetRoleByID(ctx, role.ID)
		rq.NoError(err)

		role.CreatedAt = time.Unix(role.CreatedAt.Unix(), 0)
		role.Auths = []string{auth1}
		actual.CreatedAt = time.Unix(actual.CreatedAt.Unix(), 0)

		rq.EqualValues(role, actual)
	})

	t.Run("read a deleted role", func(t *testing.T) {
		role := datastore.Role{
			RoleName: "test_role_read_a_deleted_role",
			Scopes:   []string{"scope1", "scope2"},
			Auths:    []string{auth1, auth2},
		}
		rq.NoError(store.CreateRole(ctx, &role))
		rq.NoError(store.DeleteRoleByID(ctx, role.ID))

		_, err := store.GetRoleByID(ctx, role.ID)
		rq.Equal(datastore.ErrorRoleNotExist, err)
	})

	t.Run("assign/unassign scopes", func(t *testing.T) {
		role := datastore.Role{
			RoleName: "test_role_assign_scopes",
			Scopes:   []string{"scope1", "scope2"},
			Auths:    []string{auth1, auth2},
		}
		rq.NoError(store.CreateRole(ctx, &role))

		rq.NoError(store.UpdateScopesByID(ctx, role.ID, datastore.UpdateRoleScopeOption{
			Assign:   []string{"scope3"},
			Unassign: []string{"scope1"},
		}))

		actual, err := store.GetRoleByID(ctx, role.ID)
		rq.NoError(err)
		rq.EqualValues([]string{"scope2", "scope3"}, actual.Scopes)
	})

	t.Run("assign duplicated scopes", func(t *testing.T) {
		role := datastore.Role{
			RoleName: "test_role_assign_deplicated_scopes",
			Scopes:   []string{"scope1", "scope2"},
			Auths:    []string{auth1, auth2},
		}
		rq.NoError(store.CreateRole(ctx, &role))

		rq.NoError(store.UpdateScopesByID(ctx, role.ID, datastore.UpdateRoleScopeOption{
			Assign: []string{"scope2"},
		}))

		actual, err := store.GetRoleByID(ctx, role.ID)
		rq.NoError(err)
		rq.EqualValues([]string{"scope1", "scope2"}, actual.Scopes)
	})

	t.Run("unassign non-existed scopes", func(t *testing.T) {
		role := datastore.Role{
			RoleName: "test_role_unassign_non-existed_scopes",
			Scopes:   []string{"scope1", "scope2"},
			Auths:    []string{auth1, auth2},
		}
		rq.NoError(store.CreateRole(ctx, &role))

		rq.Equal(datastore.ErrorUnassignNonExistedScopes,
			store.UpdateScopesByID(ctx, role.ID, datastore.UpdateRoleScopeOption{
				Unassign: []string{"scope3"},
			}),
		)
	})
}
