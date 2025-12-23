package config

import (
	"sync"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

var memoryRedisConnOnce sync.Once
var memoryRedisConn *redis.Client

func NewInMemoryRedis() *redis.Client {
	if memoryRedisConn == nil {
		memoryRedisConnOnce.Do(
			func() {
				memoryRedisConn = openInMemoryRedisConn()
			},
		)
	}

	return memoryRedisConn
}

func openInMemoryRedisConn() *redis.Client {
	miniRedis, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	conn := redis.NewClient(
		&redis.Options{
			Addr: miniRedis.Addr(),
		},
	)

	return conn
}
