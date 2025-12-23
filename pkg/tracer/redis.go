package tracer

import (
	redistrace "github.com/lsgndln/dd-trace-go/contrib/redis/go-redis.v9"
	"github.com/redis/go-redis/v9"
)

func RedisTracer(opts *redis.Options) redis.UniversalClient {
	return redistrace.NewClient(opts)
}
