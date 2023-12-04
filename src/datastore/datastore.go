package datastore

import (
	"context"
	"errors"
)

type Datastore interface {
	CreateAuthority(ctx context.Context, auth *Authority) error
	DeleteAuthorityByID(ctx context.Context, id int64, force bool) error
	GetAuthorityByID(ctx context.Context, id int64) (*Authority, error)

	CreateRole(ctx context.Context, role *Role) error
	DeleteRoleByID(ctx context.Context, id int64) error
	GetRoleByID(ctx context.Context, id int64) (*Role, error)
	UpdateScopesByID(ctx context.Context, id int64, op UpdateRoleScopeOption) error
}

type UpdateRoleScopeOption struct {
	Assign   []string `json:"assign,omitempty"`
	Unassign []string `json:"unassign,omitempty"`
}

var (
	ErrorAuthExist             = errors.New("authority exist")
	ErrorAuthNotExist          = errors.New("authority not exist")
	ErrorDeleteAuthWithBinding = errors.New("delete an auth with bindings")

	ErrorRoleExist    = errors.New("role exist")
	ErrorRoleNotExist = errors.New("role not exist")

	//ErrorScopesDuplicatedAssign   = errors.New("try to assign scopes which have been assigned to role")

	ErrorUnassignNonExistedScopes = errors.New("unassign non-existed scopes")
)
