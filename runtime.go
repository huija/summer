package summer

import (
	"github.com/huija/summer/conf"
	"github.com/huija/summer/container/pipeline"
	"github.com/huija/summer/logs"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

var closePipe = pipeline.NewSimpleFuncStage("", "")
var sc chan os.Signal
var rt sync.Once // avoid repeat run when test

func runtimeSetup() (err error) {
	// set the max num of running cpu thread
	runtime.GOMAXPROCS(runtime.NumCPU())
	// signal monitor & exit gracefully
	sc = make(chan os.Signal, 1)

	RegisterClose(Runtime, func() error {
		close(sc)
		logs.SugaredLogger.Info("runtime sc close...")
		return nil
	}, pipeline.ClosePipeLevel(pipeline.Enhancer))

	signal.Notify(sc,
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM,
		syscall.SIGKILL)
	go rt.Do(signalMonitor)
	logs.SugaredLogger.Debug("runtime init successfully")
	return
}

func signalMonitor() {
	for {
		for sig := range sc {
			switch sig {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM:
				logs.SugaredLogger.Warnf("Got exit signal %v.", sig)
				logs.SugaredLogger.Infof("%+v", closePipe)
				_, _ = closePipe.Exec("")
				os.Exit(0)
			case syscall.SIGKILL:
				os.Exit(1)
			case conf.SIGOK:
				logs.SugaredLogger.Infof("Got Test OK signal %v.", sig)
			case conf.SIGTODO:
				logs.SugaredLogger.Infof("Got Test TODO signal %v.", sig)
			default:
				logs.SugaredLogger.Infof("Got ignore signal %v.", sig)
			}
		}
	}
}
