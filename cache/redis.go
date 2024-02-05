package cache

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	redis *redis.Client
}

func NewRedisCache(address string, port int, password string, db int) (*RedisCache, error) {
	var redisCache RedisCache
	address = fmt.Sprintf("%s:%s", address, strconv.Itoa(port))
	redisCache.redis = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := redisCache.redis.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return &redisCache, nil
}

func (r *RedisCache) Get(ctx context.Context, key string) (*CacheValue, error) {
	if err := r.IsValid(); err != nil {
		return nil, err
	}
	v, err := r.redis.Get(ctx, key).Result()
	if err != nil || err == redis.Nil {
		return nil, err
	}
	return NewCacheValue(v)
}

func (r *RedisCache) Set(ctx context.Context, key string, value any, expiresAfter int64) error {
	if err := r.IsValid(); err != nil {
		return err
	}
	expires := time.Duration(expiresAfter * int64(time.Second))
	err := r.redis.Set(ctx, key, value, expires).Err()
	if err != nil {
		return err
	}
	return nil
}

//func (r *RedisCache) SetJson(ctx context.Context, key string, value any, expiresAfter int64) error {
//	if err := r.IsValid(); err != nil {
//		return err
//	}
//	r.redis.JSONSet(ctx, key, "$", value)
//	expires := time.Duration(expiresAfter * int64(time.Second))
//	err := r.redis.Set(ctx, key, value, expires).Err()
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	if err := r.IsValid(); err != nil {
		return err
	}
	err := r.redis.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisCache) IsValid() error {
	if r.redis == nil {
		return errors.New("redis client not declared")
	}
	return nil
}
