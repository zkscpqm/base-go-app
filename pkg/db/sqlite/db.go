package sqlite

import (
	"database/sql"
	"fmt"
	"strings"
	"unnamed/pkg/common"
	"unnamed/pkg/common/config"
	"unnamed/pkg/common/log"
	"unnamed/pkg/db"

	_ "embed"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed scripts/init.sql
var initScript string

func init() {
	db.RegisterDriver("sqlite", New)
}

type Database struct {
	dbo    *sql.DB
	logger *log.Logger
	cache  common.CacheManager
}

func New(cfg config.DbConfig, logger *log.Logger) (db.Database, error) {
	if !strings.HasSuffix(cfg.Database, ".db") {
		cfg.Database += ".db"
	}
	dbo, err := sql.Open("sqlite3", cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite DB [%s]: %v", cfg.Database, err)
	}
	return &Database{dbo: dbo, logger: logger.NewFrom("SQLITE"), cache: common.NewCacheManager()}, nil
}

func (d *Database) InitSchema() (err error) {

	if _, err = d.dbo.Exec(initScript); err != nil {
		d.logger.Debug("SQLITE INIT ERR: %v\nScript:\n%s", err, initScript)
		return fmt.Errorf("failed to execute init script: %v", err)
	}
	return
}

func (d *Database) Close() error {
	return d.dbo.Close()
}

func (d *Database) InSession(f func(tx *db.Session) error) error {
	tx, err := d.dbo.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	if err = f(&db.Session{Tx: tx}); err != nil {
		return fmt.Errorf("transaction failed: %v; rollback error: %v", err, tx.Rollback())
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	d.logger.Debug("Transaction committed!")
	return nil
}
