package summer

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/huija/summer/conf"
	"github.com/huija/summer/container/pipeline"
	"github.com/huija/summer/logs"
	"github.com/huija/summer/srv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"
)

func testInit() (err error) {
	// default use empty `conf/config.yaml`
	//*confPath = "./conf/config_template.yaml"

	gin.SetMode(gin.TestMode)

	s := CustomizeStage(GinRun, func() error { return nil }, pipeline.Runner)
	if s == nil {
		err = errors.New("test init: customize stage failed")
		return
	}
	return Bloom()
}

func srvPing(t *testing.T) {
	defer func() {
		if t.Failed() {
			sc <- syscall.SIGKILL
		} else {
			sc <- conf.SIGOK
			err := srv.Server.Shutdown(context.Background())
			require.Equal(t, nil, err)
		}
	}()
	time.Sleep(1 * time.Second) // wait server start
	r, err := http.Get("http://localhost:9090/ping")
	require.Equal(t, nil, err)
	assert.Equal(t, http.StatusOK, r.StatusCode)
	body, err := ioutil.ReadAll(r.Body)
	require.Equal(t, nil, err)
	assert.Equal(t, `{"ping":"pong"}`, string(body))
}

func TestMain(m *testing.M) {
	err := testInit()
	if err != nil {
		logs.SugaredLogger.Debug(err.Error())
		os.Exit(1)
	}
	code := m.Run()
	sc <- syscall.SIGPIPE // ignore
	if code == 0 {
		//sc <- syscall.SIGTERM
	} else {
		sc <- syscall.SIGKILL
	}
	time.Sleep(1 * time.Millisecond) // wait for runtime signal monitor
}
