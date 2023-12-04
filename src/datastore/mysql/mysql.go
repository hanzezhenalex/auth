package mysql

import (
	"context"
	"fmt"

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

	if src.IsDebugMode() {
		store.engine.ShowSQL(true)
	}

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
		new(datastore.RoleBinding),
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
		if _, err := store.engine.
			Table(table).
			Where("1=1").
			Unscoped().
			Delete(); err != nil {
			name := table.(interface{ TableName() string }).TableName()
			return fmt.Errorf("fail to delete data in table %s, err=%w", name, err)
		}
	}
	return nil
}

func (store *mysqlDatastore) transaction(ctx context.Context, fn func(*xorm.Session) error) error {
	session := store.engine.NewSession().Context(ctx)
	defer func() { _ = session.Close() }()
	rollback := func() { _ = session.Rollback() }

	if err := session.Begin(); err != nil {
		return fmt.Errorf("fail to start session: %w", err)
	}

	if err := fn(session); err != nil {
		rollback()
		return err
	}

	if err := session.Commit(); err != nil {
		return fmt.Errorf("fail to commit session: %w", err)
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
func (store *mysqlDatastore) DeleteAuthorityByID(ctx context.Context, id int64, force bool) error {
	return store.transaction(ctx, func(session *xorm.Session) error {
		if !force {
			cnt, err := session.
				Where("auth_id=?", id).
				Count(new(datastore.RoleBinding))
			if err != nil {
				return fmt.Errorf("fail to count role bindings, %w", err)
			}
			if cnt != 0 {
				return datastore.ErrorDeleteAuthWithBinding
			}
		}

		n, err := session.
			Table(new(datastore.Authority)).
			Where("id=?", id).
			Delete()

		if err != nil {
			return err
		} else if n == 0 {
			return datastore.ErrorAuthNotExist
		} else {
			return nil
		}
	})

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

/*
	Role
*/

func (store *mysqlDatastore) CreateRole(ctx context.Context, role *datastore.Role) error {
	return store.transaction(ctx, func(session *xorm.Session) error {
		// step 1: insert role
		if _, err := session.Insert(role); err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok {
				if mysqlErr.Number == duplicatedOnPrimaryKey {
					return datastore.ErrorRoleExist
				}
			}
			return fmt.Errorf("fail to insert role: %w", err)
		}

		// step 2: fetch authorities
		var auths []datastore.Authority
		if err := session.In("authority_name", role.Auths).
			Table(new(datastore.Authority)).
			Find(&auths); err != nil {
			return fmt.Errorf("fail to fetch auths: %w", err)
		}

		// step 3: check if all needed authorities exist
		if len(auths) != len(role.Auths) {
			return datastore.ErrorAuthNotExist
		}

		// step 4: insert role-auth relationship if needed
		if len(auths) > 0 {
			rbs := make([]datastore.RoleBinding, 0, len(auths))
			for _, auth := range auths {
				rbs = append(rbs, datastore.RoleBinding{
					RoleID:   role.ID,
					AuthID:   auth.ID,
					AuthName: auth.AuthName,
				})
			}
			if _, err := session.InsertMulti(&rbs); err != nil {
				return fmt.Errorf("fail to insert rbs, %w", err)
			}
		}
		return nil
	})
}

func (store *mysqlDatastore) DeleteRoleByID(ctx context.Context, id int64) error {
	return store.transaction(ctx, func(session *xorm.Session) error {
		// step 1: delete role-auth bindings
		if _, err := session.
			Table(new(datastore.RoleBinding)).
			Where("role_id=?", id).
			Delete(); err != nil {
			return fmt.Errorf("fail to delete role bindings, %w", err)
		}

		// step 2: delete role
		n, err := session.
			Table(new(datastore.Role)).
			Where("id=?", id).
			Delete()

		if err != nil {
			return fmt.Errorf("fail to delete role, %w", err)
		}
		if n == 0 {
			return datastore.ErrorRoleNotExist
		}
		return nil
	})
}

func (store *mysqlDatastore) GetRoleByID(ctx context.Context, id int64) (*datastore.Role, error) {
	var role datastore.Role
	err := store.transaction(ctx, func(session *xorm.Session) error {
		if ok, err := session.
			ID(id).
			Get(&role); err != nil {
			return fmt.Errorf("fail to get role: %w", err)
		} else if !ok {
			return datastore.ErrorRoleNotExist
		}

		results, err := session.QueryString(getActiveRoleAuthNames(id))
		if err != nil {
			return fmt.Errorf("fail to get role binding, %w", err)
		}

		for _, rb := range results {
			role.Auths = append(role.Auths, rb["auth_name"])
		}
		return nil
	})
	return &role, err
}

func (store *mysqlDatastore) UpdateScopesByID(ctx context.Context, id int64, op datastore.UpdateRoleScopeOption) error {
	return store.transaction(ctx, func(session *xorm.Session) error {
		var role datastore.Role
		if ok, err := session.
			ForUpdate().
			ID(id).
			Get(&role); err != nil {
			return fmt.Errorf("fail to get role %d, %w", id, err)
		} else if !ok {
			return datastore.ErrorRoleNotExist
		}

		scopesAppend, duplicated := src.SliceAppend(role.Scopes, op.Assign)
		if len(duplicated) > 0 {
			//return datastore.ErrorScopesDuplicatedAssign
		}

		scopesRemoved, nonExisted := src.SliceRemove(scopesAppend, op.Unassign)
		if len(nonExisted) > 0 {
			return datastore.ErrorUnassignNonExistedScopes
		}

		src.SortSliceAsc(scopesRemoved)
		role.Scopes = scopesRemoved
		if _, err := session.
			ID(id).
			Cols("scopes").
			Update(&role); err != nil {
			return err
		}
		return nil
	})
}
