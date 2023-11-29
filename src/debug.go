package src

import (
	"fmt"
	"os"
	"sync"
)

const (
	envMode   = "AUTH_RUNNING_MODE"
	envUserID = "USER_ID"
)

/*
	Mode
*/

const (
	DebugMode = "debug"
	//ProductionMode = "production"
)

func EnableDebugMode() {
	_ = os.Setenv(envMode, DebugMode)
}

func IsDebugMode() bool {
	mode := os.Getenv(envMode)
	return mode == DebugMode
}

/*
	User ID
*/

var (
	userID string
	once   sync.Once
)

func GetUserID() string {
	once.Do(func() {
		id := os.Getenv(envUserID)
		if id == "" {
			raw, err := GenerateSecureRandomString(10)
			if err != nil {
				panic(err)
			}
			id = raw
		}
		userID = id
	})
	return userID
}

func WithDebugSuffix(name string) string {
	if IsDebugMode() {
		return fmt.Sprintf("%s-%s", name, GetUserID())
	}
	return name
}
