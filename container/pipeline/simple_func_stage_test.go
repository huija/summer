package pipeline

import (
	"github.com/stretchr/testify/require"
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

func initTest(t *testing.T) {
	root = NewSimpleFuncStage("", "")
	p := root.AddStage("", config,
		NewSimpleFuncStage(config, "").SetFunc(func() error {
			t.Log(config)
			return nil
		}))
	require.NotEqual(t, nil, p)

	p = root.AddStage(webFrameworkTrace, webFrameworkTraceChild,
		NewSimpleFuncStage(webFrameworkTraceChild, "").SetFunc(func() error {
			t.Log(webFrameworkTraceChild)
			return nil
		}))
	require.NotEqual(t, nil, p)

	p = root.AddStage(webFramework, webFrameworkTrace,
		NewSimpleFuncStage(webFrameworkTrace, "").SetFunc(func() error {
			t.Log(webFrameworkTrace)
			return nil
		}))
	require.NotEqual(t, nil, p)

	p = root.CustomizeStage(webFramework,
		NewSimpleFuncStage(webFramework, "").SetFunc(func() error {
			t.Log(webFramework)
			return nil
		}).SetPriority(Runner))
	require.NotEqual(t, nil, p)

	p = root.AddStage(webFramework, webFrameworkRun,
		NewSimpleFuncStage(webFrameworkRun, "").SetFunc(func() error {
			t.Log(webFrameworkRun)
			return nil
		}).SetPriority(Runner))
	require.NotEqual(t, nil, p)
}

func TestSFSExec(t *testing.T) {
	initTest(t)

	// aim: config->wf->wfTrace->wfRun
	root.Exec(nil)
}

func TestSFSCustomize(t *testing.T) {
	initTest(t)

	p := root.CustomizeStage(webFrameworkRun,
		NewSimpleFuncStage(webFrameworkRun, "").SetFunc(func() error {
			t.Log(webFrameworkRun + "-" + webFrameworkRun)
			return nil
		}).SetPriority(Initiator))
	require.NotEqual(t, nil, p)

	// aim: config->wf->wfRun->wfTrace
	root.Exec(nil)
}

func TestSFSRemove(t *testing.T) {
	initTest(t)

	p := root.RemoveStage(webFrameworkTraceChild)
	require.NotEqual(t, nil, p)

	// aim: config->wf->wfTrace->wfRun
	t.Log(root.String())
	root.Exec(nil)
}
