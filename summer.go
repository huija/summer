package summer

import (
	"github.com/huija/summer/container/pipeline"
	"github.com/huija/summer/logs"
	"github.com/huija/summer/utils"
	"strings"
)

const (
	Split = "."

	// init core
	Config  = "config"
	Log     = "log"
	Runtime = "runtime"

	// debug middleware
	DBs = "dbs"

	// runner
	Gin = "gin"
)

var rootPipe = pipeline.NewSimpleFuncStage("", "")

var removed []string

// GetRootPipe get root pipe
func GetRootPipe() *pipeline.SimpleFuncStage {
	return rootPipe
}

func init() {
	AddStage(Config, yamlConfig, pipeline.Initiator)
	AddStage(Log, logger, pipeline.Debugger)
	AddStage(Runtime, runtimeSetup, pipeline.Enhancer)
	AddStage(DBs, databases)
	AddStage(Gin, ginEngine, pipeline.Runner)
}

// AddStage to pipe
// for user, pipeline.Feature is default level instead of pipeline.Debugger(zero)
func AddStage(key string, f func() error, p ...pipeline.PipeLevel) pipeline.PipeLine {
	if utils.ArrayIndexOf(removed, key) != -1 {
		return nil
	}

	index := strings.LastIndex(key, Split)
	if index == -1 {
		index = 0
	}

	// default priority Feature
	return rootPipe.AddStage(key[:index], key,
		pipeline.NewSimpleFuncStage(key, Split).
			SetFunc(f).SetPriority(append(p, pipeline.Feature)[0]))
}

// CustomizeStage in pipe
// customize before add will avoid repeat add
func CustomizeStage(key string, f func() error, p pipeline.PipeLevel) pipeline.PipeLine {
	if utils.ArrayIndexOf(removed, key) != -1 {
		return nil
	}

	if stage := rootPipe.GetStage(key); stage != nil {
		return stage.CustomizeStage(key,
			pipeline.NewSimpleFuncStage(key, Split).
				SetFunc(f).SetPriority(p))
	}

	return AddStage(key, f, p)
}

// RemoveStage in pipe
// remove before add will avoid add
func RemoveStage(key string) pipeline.PipeLine {
	removed = append(removed, key)
	return rootPipe.RemoveStage(key)
}

// GetStage in pipe
func GetStage(key string) pipeline.PipeLine {
	return rootPipe.GetStage(key)
}

// RegisterClose register close functions
func RegisterClose(key string, fn func() error, p ...pipeline.PipeLevel) {
	index := strings.LastIndex(key, Split)
	if index == -1 {
		index = 0
	}

	// default priority Feature
	level := append(p, pipeline.ClosePipeLevel(pipeline.Feature))[0]
	closePipe.AddStage(key[:index], key,
		pipeline.NewSimpleFuncStage(key, Split).
			SetFunc(fn).SetPriority(level))
}

// Bloom like summer flowers
func Bloom() error {
	logs.SugaredLogger.Infof("pipe: %+v", rootPipe)

	_, err := rootPipe.Exec("")
	return err
}
