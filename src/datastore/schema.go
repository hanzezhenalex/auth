package datastore

import (
	"encoding/json"
	"fmt"
	"github.com/hanzezhenalex/auth/src"
	"strings"
	"time"
)

type User struct {
	ID        int64     `xorm:"'id' pk autoincr"`
	Username  string    `xorm:"'user_name' not null unique(is_delete)"`
	Password  string    `xorm:"'password' not null"`
	Reserve   string    `xorm:"'reserve'"`
	CreatedAt time.Time `xorm:"created"`
	DeletedAt int64     `xorm:"'deleted_at' unique(is_delete) default(0)"`
}

func (user User) TableName() string {
	return src.WithDebugSuffix("user")
}

type Scopes []string

const delimiter = ";"

func (s Scopes) MarshalJSON() ([]byte, error) {
	result := ""
	for _, scope := range s {
		result += delimiter + scope
	}
	return json.Marshal(result[1:])
}

func (s *Scopes) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("unable to unmarshal, err=%s", err)
	}
	scopes := strings.Split(raw, delimiter)
	*s = scopes
	return nil
}

type Role struct {
	ID        int64     `xorm:"'id' pk autoincr"`
	RoleName  string    `xorm:"'role_name' unique(is_delete)"`
	CreatedBy string    `xorm:"'created_by'"`
	Scopes    Scopes    `xorm:"'scopes'"`
	CreatedAt time.Time `xorm:"created"`
	DeletedAt int64     `xorm:"'deleted_at' unique(is_delete) default(0)"`
}

func (role Role) TableName() string {
	return src.WithDebugSuffix("role")
}

type Authority struct {
	ID        int64     `xorm:"'id' pk autoincr"`
	AuthName  string    `xorm:"'authority_name' not null unique(is_delete)"`
	CreatedBy string    `xorm:"'created_by'"`
	CreatedAt time.Time `xorm:"created"`
	DeletedAt int64     `xorm:"'deleted_at' unique(is_delete) default(0)"`
}

func (auth Authority) TableName() string {
	return src.WithDebugSuffix("authority")
}

type RoleBinding struct {
}

type AuthBinding struct {
}
