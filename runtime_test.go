package summer

import (
	"github.com/huija/summer/conf"
	"testing"
	"time"
)

func TestInitRuntime(t *testing.T) {
	sc <- conf.SIGTODO
	// TODO single test runtime
	//sc <- syscall.SIGTERM // then refer to log to see details
	time.Sleep(1 * time.Second)
}
