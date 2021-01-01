package logs

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestDefaults(t *testing.T) {
	defaults, err := Defaults(nil)
	require.Equal(t, nil, err)
	a, err := yaml.Marshal(defaults)
	require.Equal(t, nil, err)
	t.Log("\n", string(a))

	logs, err := Defaults([]*Log{{
		Type: File,
	}})
	require.Equal(t, nil, err)
	b, err := yaml.Marshal(logs)
	require.Equal(t, nil, err)
	t.Log("\n", string(b))

	logs, err = Defaults([]*Log{{
		Type: "",
	}})
	require.Equal(t, nil, err)
	b, err = yaml.Marshal(logs)
	require.Equal(t, nil, err)
	t.Log("\n", string(b))
}

func TestInnerLog(t *testing.T) {
	SugaredLogger.Debug("1")
	SugaredLogger.Debugf("%d", 1)
	SugaredLogger.Info("2")
	SugaredLogger.Infof("%d", 2)
	SugaredLogger.Warn("3")
	SugaredLogger.Warnf("%d", 3)
	SugaredLogger.Error("4")
	SugaredLogger.Errorf("%d", 4)
	//SugaredLogger.Panic("5")
	//SugaredLogger.Panicf("%d", 5)
	//SugaredLogger.Fatal("6")
	//SugaredLogger.Fatalf("%d", 6)
}
