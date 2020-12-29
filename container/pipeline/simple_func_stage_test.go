package pipeline

import (
	"errors"
	"github.com/huija/summer/logs"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	config                 = "config"
	webFramework           = "wf"
	webFrameworkRun        = "wfRun"
	webFrameworkTrace      = "wfTrace"
	webFrameworkTraceChild = "wfTraceChild"
)

var root *SimpleFuncStage

func initTest() error {
	root = NewSimpleFuncStage("", "")
	p := root.AddStage("", config,
		NewSimpleFuncStage(config, "").SetFunc(func() error {
			logs.SugaredLogger.Debug(config)
			return nil
		}))
	if p == nil {
		return errors.New("test init: pipeline is nil")
	}

	p = root.AddStage(webFrameworkTrace, webFrameworkTraceChild,
		NewSimpleFuncStage(webFrameworkTraceChild, "").SetFunc(func() error {
			logs.SugaredLogger.Debug(webFrameworkTraceChild)
			return nil
		}))
	if p == nil {
		return errors.New("test init: pipeline is nil")
	}

	p = root.AddStage(webFramework, webFrameworkTrace,
		NewSimpleFuncStage(webFrameworkTrace, "").SetFunc(func() error {
			logs.SugaredLogger.Debug(webFrameworkTrace)
			return nil
		}))
	if p == nil {
		return errors.New("test init: pipeline is nil")
	}

	p = root.CustomizeStage(webFramework,
		NewSimpleFuncStage(webFramework, "").SetFunc(func() error {
			logs.SugaredLogger.Debug(webFramework)
			return nil
		}).SetPriority(Runner))
	if p == nil {
		return errors.New("test init: pipeline is nil")
	}

	p = root.AddStage(webFramework, webFrameworkRun,
		NewSimpleFuncStage(webFrameworkRun, "").SetFunc(func() error {
			logs.SugaredLogger.Debug(webFrameworkRun)
			return nil
		}).SetPriority(Runner))
	if p == nil {
		return errors.New("test init: pipeline is nil")
	}
	return nil
}

func TestSimpleFuncStage_CustomizeStage(t *testing.T) {
	p := root.CustomizeStage(webFrameworkRun,
		NewSimpleFuncStage(webFrameworkRun, "").SetFunc(func() error {
			logs.SugaredLogger.Debug(webFrameworkRun + "-" + webFrameworkRun)
			return nil
		}).SetPriority(Initiator))
	require.NotEqual(t, nil, p)

	// aim: config->wf->wfRun->wfTrace
	t.Log(root.String())
}

func TestSimpleFuncStage_RemoveStage(t *testing.T) {
	p := root.RemoveStage(webFrameworkTraceChild)
	require.NotEqual(t, nil, p)

	// aim: config->wf->wfTrace->wfRun
	t.Log(root.String())
}

func TestSimpleFuncStage_GetStage(t *testing.T) {
	p := root.GetStage(webFrameworkTraceChild)
	require.NotEqual(t, nil, p)
}

func TestSimpleFuncStage_Exec(t *testing.T) {
	// aim: config->wf->wfTrace->wfRun
	t.Log(root.String())

	root.Exec(nil)
}

func TestClosePipeLevel(t *testing.T) {
	require.Equal(t, ClosePipeLevel(Initiator), Runner)
	require.Equal(t, ClosePipeLevel(Initiator+1), Runner)
	require.Equal(t, ClosePipeLevel(Runner), Initiator+1)
	require.Equal(t, ClosePipeLevel(Debugger), Debugger)
}

func TestMain(m *testing.M) {
	err := initTest()
	if err != nil {
		os.Exit(1)
	}
	code := m.Run()
	if code != 0 {
		os.Exit(code)
	}
}
