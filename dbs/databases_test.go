package dbs

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDefaults(t *testing.T) {
	var (
		dbs = &DBs{
			Redis: &Redis{
				//Schema: "Not Equal",
			},
			MySQL: &MySQL{},
			Mongo: &Mongo{}}
		err      error
		defaults = &DBs{
			MySQL: defaultsMySQL(),
			Mongo: defaultsMongo(),
			Redis: defaultsRedis(),
		}
	)
	dbs, err = Defaults(dbs)
	require.Equal(t, nil, err)
	require.Equal(t, defaults, dbs)
}
