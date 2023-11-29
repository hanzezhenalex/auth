package datastore

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type User struct {
	Username  string    `xorm:"user_name not null"`
	Password  string    `xorm:"password not null"`
	ID        string    `xorm:"id  pk"`
	Binding   string    `xorm:"binding"`
	CreatedAt time.Time `xorm:"created"`
	DeletedAt time.Time `xorm:"deleted_at"`
}

func (user User) TableName() string {
	return "user"
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
	RoleName  string    `xorm:"role_name pk"`
	CreatedBy string    `xorm:"created_by"`
	Scopes    Scopes    `xorm:"scopes"`
	CreatedAt time.Time `xorm:"created"`
	DeletedAt time.Time `xorm:"deleted_at"`
}

func (role Role) TableName() string {
	return "role"
}

type Authority struct {
	AuthName  string    `xorm:"authority_name pk"`
	CreatedBy string    `xorm:"created_by"`
	CreatedAt time.Time `xorm:"created"`
	DeletedAt time.Time `xorm:"deleted_at"`
}

func (auth Authority) TableName() string {
	return "authority"
}

type RoleBinding struct {
}

type AuthBinding struct {
}
