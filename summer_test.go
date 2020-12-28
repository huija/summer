package summer

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/huija/summer/conf"
	"github.com/huija/summer/container/pipeline"
	"github.com/huija/summer/srv"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"syscall"
	"testing"
	"time"
)

func testInit(t *testing.T) {
	// TODO default use empty `conf/config.yaml`, delete next line
	*confPath = "./conf/config_template.yaml"

	var err error
	gin.SetMode(gin.ReleaseMode)
	s := CustomizeStage(GinRun, func() error { return nil }, pipeline.Runner)
	require.NotEqual(t, nil, s)
	err = Bloom()
	require.Equal(t, nil, err)
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
	require.Equal(t, http.StatusOK, r.StatusCode)
	body, err := ioutil.ReadAll(r.Body)
	require.Equal(t, nil, err)
	require.Equal(t, `{"ping":"pong"}`, string(body))
}

func TestMain(m *testing.M) {
	testInit(new(testing.T))
	code := m.Run()
	sc <- syscall.SIGPIPE // ignore
	if code == 0 {
		//sc <- syscall.SIGTERM
	} else {
		sc <- syscall.SIGKILL
	}
	time.Sleep(1 * time.Millisecond) // wait for runtime signal monitor
}
