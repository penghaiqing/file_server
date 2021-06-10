package redis

import (
	"github.com/garyburd/redigo/redis"
)

var (
	pool      *redis.Pool
	redisHost = "127.0.0.1:6379"
	redisPass = "root"
)

// newRedisPool: 创建redis 连接池
func newRedisPool() *redis.Pool {
	return &redis.Pool{}
}
