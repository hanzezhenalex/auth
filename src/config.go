package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	defaultMaxIdleConns = 5
	defaultMaxOpenConns = 10
)

type DbConfig struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	Database     string `json:"database"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	MaxIdleConns int    `json:"max_idle_conns,omitempty"`
	MaxOpenConns int    `json:"max_open_conns,omitempty"`
}

func NewDbConfig() DbConfig {
	return DbConfig{
		MaxIdleConns: defaultMaxIdleConns,
		MaxOpenConns: defaultMaxOpenConns,
	}
}

func (dbCfg DbConfig) Dns() string {
	// "username:password@tcp(host:post)/dbname"
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		dbCfg.Username, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Database)
	return dns
}

type Config struct {
	DbConfig
}

func NewConfigFromFile(path string) (Config, error) {
	var cfg = Config{
		NewDbConfig(),
	}
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("fail to read config file, %w", err)
	}

	if err := json.NewDecoder(bytes.NewBuffer(raw)).Decode(&cfg); err != nil {
		return cfg, fmt.Errorf("fail to decode config file, %w", err)
	}
	return cfg, err
}
