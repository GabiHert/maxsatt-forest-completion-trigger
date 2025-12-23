package redishelper

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type loggerAdapter interface {
	Debug(ctx context.Context, message string, metadata ...any)
}

type RedisHelper interface {
	Put(ctx context.Context, id string, object any, optionalDuration ...time.Duration) error
	Get(ctx context.Context, id string, object any) error
	Delete(ctx context.Context, id string) error
	Clear(ctx context.Context, prefix string) error
}

type redisHelper struct {
	redis        redis.UniversalClient
	logger       loggerAdapter
	globalPrefix string
}

func Redis(redis redis.UniversalClient, globalPrefix string, logger loggerAdapter) RedisHelper {
	return &redisHelper{
		redis:        redis,
		globalPrefix: globalPrefix,
		logger:       logger,
	}
}

func (r *redisHelper) Put(ctx context.Context, id string, object any, optionalDuration ...time.Duration) error {
	r.logger.Debug(ctx, "Started", id)

	var value []byte
	switch object.(type) {
	case string:
		value = []byte(object.(string))
	case []byte:
		value = object.([]byte)
	default:
		valueBytes, err := json.Marshal(object)
		if err != nil {
			return err
		}
		value = valueBytes
	}

	var duration time.Duration
	if len(optionalDuration) > 0 {
		duration = optionalDuration[0]
	} else {
		duration = time.Duration(10) * time.Hour
	}

	err := r.redis.SetEx(ctx, r.globalPrefix+id, value, duration).Err()
	if err != nil {
		return err
	}

	r.logger.Debug(ctx, "Finished", id)
	return nil
}

func (r *redisHelper) Get(ctx context.Context, id string, object any) error {
	r.logger.Debug(ctx, "Started", id)

	result, err := r.redis.Get(ctx, r.globalPrefix+id).Result()
	if err != nil {
		return err
	}

	switch object.(type) {
	case *string:
		*(object.(*string)) = result
	case *[]byte:
		*(object.(*[]byte)) = []byte(result)
	default:
		if err := json.Unmarshal([]byte(result), &object); err != nil {
			return err
		}
	}

	r.logger.Debug(ctx, "Finished", id)
	return nil
}

func (r *redisHelper) Delete(ctx context.Context, id string) error {
	r.logger.Debug(ctx, "Started", id)

	err := r.redis.Del(ctx, r.globalPrefix+id).Err()
	if err != nil {
		return err
	}

	r.logger.Debug(ctx, "Finished", id)
	return nil
}

func (r *redisHelper) Clear(ctx context.Context, prefix string) error {
	r.logger.Debug(ctx, "Started", prefix)

	var cursor uint64
	var keys []string
	for {
		var tmpKeys []string
		var err error
		tmpKeys, cursor, err = r.redis.Scan(ctx, cursor, r.globalPrefix+prefix, 100).Result()
		if err != nil {
			return err
		}
		keys = append(keys, tmpKeys...)
		if cursor == 0 {
			break
		}
	}

	for _, key := range keys {
		err := r.redis.Del(ctx, key).Err()
		if err != nil {
			return err
		}
	}

	r.logger.Debug(ctx, "Finished", prefix)
	return nil
}

// NoOpRedisHelper returns a RedisHelper that does nothing (for when Redis is not configured)
func NoOpRedisHelper() RedisHelper {
	return &noOpRedisHelper{}
}

type noOpRedisHelper struct{}

func (n *noOpRedisHelper) Put(ctx context.Context, id string, object any, optionalDuration ...time.Duration) error {
	return nil
}

func (n *noOpRedisHelper) Get(ctx context.Context, id string, object any) error {
	return errors.New("redis not configured")
}

func (n *noOpRedisHelper) Delete(ctx context.Context, id string) error {
	return nil
}

func (n *noOpRedisHelper) Clear(ctx context.Context, prefix string) error {
	return nil
}
