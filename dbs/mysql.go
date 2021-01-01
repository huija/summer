package dbs

import (
	"database/sql"
)

// MySQL mysql setup
type MySQL struct {
	Schema      string `json:",omitempty"`
	MaxPoolSize int    `json:",omitempty"`
	MinPoolSize int    `json:",omitempty"`
}

var MysqlDB *sql.DB

func defaultsMySQL() *MySQL {
	return &MySQL{
		Schema:      "root:123456@tcp(127.0.0.1:3306)/summer?charset=utf8mb4&parseTime=True&loc=Local",
		MaxPoolSize: 20,
		MinPoolSize: 5,
	}
}
