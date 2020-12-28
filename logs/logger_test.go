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
}
