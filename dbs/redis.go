package dbs

import "github.com/go-redis/redis/v8"

// Redis redis setup
type Redis struct {
	Addrs       []string
	MaxPoolSize int
	MinPoolSize int
	Username    string
	Password    string
	DB          int
}

var RedisDB redis.UniversalClient
