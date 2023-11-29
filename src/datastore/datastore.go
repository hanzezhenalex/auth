package datastore

import (
	"context"
	"errors"
)

type Datastore interface {
	CreateAuthority(ctx context.Context, auth *Authority) error
	DeleteAuthorityByID(ctx context.Context, id int64) error
}

var (
	ErrorAuthExist    = errors.New("authority exist")
	ErrorAuthNotExist = errors.New("authority not exist")
)
