package conf

import (
	"github.com/huija/summer/dbs"
	"github.com/huija/summer/logs"
	"github.com/huija/summer/srv"
	"syscall"
)

const (
	SIGOK   = syscall.Signal(0xa)
	SIGTODO = syscall.Signal(0xc)
)

var Config YamlConfig

type YamlConfig struct {
	Srv  *srv.Srv    `yaml:",omitempty"`
	Logs []*logs.Log `yaml:",omitempty"`
	DBs  *dbs.DBs    `yaml:",omitempty"`
}
