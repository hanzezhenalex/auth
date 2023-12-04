package mysql

import (
	"github.com/hanzezhenalex/auth/src/datastore"

	"xorm.io/builder"
)

func getActiveRoleAuthNames(id int64) *builder.Builder {
	return builder.
		Select("rbs.auth_name").
		From(
			builder.
				Select("auth_id", "auth_name").
				From(new(datastore.RoleBinding).TableName()).
				Where(builder.Eq{"deleted_at": 0}).
				And(builder.Eq{"role_id": id}),
			"rbs").
		LeftJoin(
			builder.
				Select("id").
				From(new(datastore.Authority).TableName()).
				Where(builder.Eq{"deleted_at": 0}),
			"id=rbs.auth_id",
			"auth").
		Where(builder.NotNull{"auth.id"})
}
