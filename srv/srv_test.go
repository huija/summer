package srv

import (
	"github.com/huija/summer/utils"
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
	err = utils.MergeStructByMarshal(srv, defaults)
	require.Equal(t, nil, err)
	marshal, err := yaml.Marshal(srv)
	require.Equal(t, nil, err)
	t.Log("\n", string(marshal))
}
