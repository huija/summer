package srv

import (
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestDefaults(t *testing.T) {
	var (
		srv      = &Srv{Host: "github.com"}
		err      error
		defaults *Srv
	)
	defaults, err = Defaults(nil)
	require.Equal(t, nil, err)
	srv, err = Defaults(srv)
	require.Equal(t, nil, err)
	require.NotEqual(t, defaults, srv)
	defaults.Host = srv.Host
	require.Equal(t, defaults, srv)

	marshal, err := yaml.Marshal(srv)
	require.Equal(t, nil, err)
	t.Log("\n", string(marshal))
}
