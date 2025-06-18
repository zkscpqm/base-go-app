package db

import "database/sql"

type Session struct {
	Tx *sql.Tx
}
