package dbs

import "github.com/go-redis/redis/v8"

// Redis redis setup
type Redis struct {
	Addrs       []string `json:",omitempty"`
	MaxPoolSize int      `json:",omitempty"`
	MinPoolSize int      `json:",omitempty"`

	// zero val is valid
	Username string
	Password string
	DB       int
}

var RedisDB redis.UniversalClient

func defaultsRedis() *Redis {
	return &Redis{
		Addrs:       []string{"127.0.0.1:6379"},
		MaxPoolSize: 20,
		MinPoolSize: 5,
	}
}
