package summer

import (
	"context"
	"github.com/huija/summer/conf"
	"github.com/huija/summer/dbs"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInitMongoDB(t *testing.T) {
	if conf.Config.DBs != nil && conf.Config.DBs.Mongo != nil {
		require.Equal(t, nil, dbs.MongoDB.Ping(context.Background(), nil))
	}
}

func TestInitMySQL(t *testing.T) {
	if conf.Config.DBs != nil && conf.Config.DBs.MySQL != nil {
		require.Equal(t, nil, dbs.MysqlDB.Ping())
	}
}

func TestInitRedis(t *testing.T) {
	if conf.Config.DBs != nil && conf.Config.DBs.Redis != nil {
		require.Equal(t, nil, dbs.RedisDB.Ping(context.Background()).Err())
	}
}
