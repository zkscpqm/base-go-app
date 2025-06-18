package db

import (
	"fmt"
	"unnamed/pkg/common/config"
	"unnamed/pkg/common/log"
)

type Driver func(cfg config.DbConfig, logger *log.Logger) (Database, error)

var drivers = make(map[string]Driver)

func RegisterDriver(name string, driver Driver) {
	drivers[name] = driver
}

type Database interface {
	InitSchema() error
	Close() error
	InSession(f func(tx *Session) error) error
}

func New(cfg config.DbConfig, logger *log.Logger) (Database, error) {
	driver, ok := drivers[cfg.Engine]
	if !ok {
		return nil, fmt.Errorf("unknown database engine: %s", cfg.Engine)
	}
	return driver(cfg, logger)
}
