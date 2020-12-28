package summer

import (
	"github.com/huija/summer/logs"
	"go.uber.org/zap"
	"sync"
	"testing"
)

func TestLogPrint(t *testing.T) {
	logs.SugaredLogger.Debug("-1")
	logs.SugaredLogger.Info("0")
	logs.SugaredLogger.Warn("1")
	logs.SugaredLogger.Error("2")
	//logs.SugaredLogger.Panic("3")
	//logs.SugaredLogger.Fatal("4")
}

func TestLogPrintf(t *testing.T) {
	logs.SugaredLogger.Debugf("%d", -1)
	logs.SugaredLogger.Infof("%d", 0)
	logs.SugaredLogger.Warnf("%d", 1)
	logs.SugaredLogger.Errorf("%d", 2)
	//logs.SugaredLogger.Panicf("%d", 3)
	//logs.SugaredLogger.Fatalf("%d", 4)
}

func TestLogConcurrent(t *testing.T) {
	var startingGun = make(chan struct{})
	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			select {
			case <-startingGun:
				logs.SugaredLogger.Debug(idx)
			}
		}(i)
	}
	close(startingGun)
	wg.Wait()
}

func TestZapLogPrint(t *testing.T) {
	ZapLogger.Debug("", zap.String("num", "-1"))
	ZapLogger.Info("", zap.String("num", "0"))
	ZapLogger.Warn("", zap.String("num", "1"))
	ZapLogger.Error("", zap.String("num", "2"))
	//ZapLogger.Panic("", zap.String("num", "3"))
	//ZapLogger.Fatal("", zap.String("num", "4"))
}
