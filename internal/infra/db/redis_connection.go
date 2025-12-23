package config

import (
	"context"
	"crypto/tls"
	"sync"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/infra/properties"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/tracer"

	"github.com/redis/go-redis/v9"
)

var redisConnOnce sync.Once
var redisConn redis.UniversalClient

func Redis() redis.UniversalClient {
	if redisConn == nil {
		redisConnOnce.Do(
			func() {
				redisConn = openRedisConn()
			},
		)
	}

	return redisConn
}

func openRedisConn() redis.UniversalClient {
	opts := &redis.Options{
		Addr:     properties.Properties().Redis.Host + ":" + properties.Properties().Redis.Port,
		Password: properties.Properties().Redis.Password,
		Username: properties.Properties().Redis.Username,
	}

	client := tracer.RedisTracer(opts)

	if _, err := client.Ping(context.TODO()).Result(); err != nil {
		opts.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		client = tracer.RedisTracer(opts)
	}

	return client
}
