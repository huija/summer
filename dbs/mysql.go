package dbs

import (
	"database/sql"
)

// MySQL mysql setup
type MySQL struct {
	Schema      string
	MaxPoolSize int
	MinPoolSize int
}

var MysqlDB *sql.DB
