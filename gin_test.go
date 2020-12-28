package summer

import (
	"github.com/gin-gonic/gin"
	"github.com/huija/summer/conf"
	"github.com/huija/summer/srv"
	"net/http"
	"testing"
)

func TestGinRun(t *testing.T) {
	t.Run("ginRun", func(t *testing.T) {
		if conf.Config.Srv != nil && conf.Config.Srv.FrameWork == srv.GIN {
			GinSrv.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"ping": "pong",
				})
			})
			go ginRun()
		}
	})

	t.Run("ginPing", srvPing)
}
