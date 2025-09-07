package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *goredis.Client
}

func NewRedisClient(addr, pass string, db int) *RedisClient {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})
	return &RedisClient{Client: rdb}
}

func (r *RedisClient) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

func (r *RedisClient) Set(ctx context.Context, key string, value string, ttlSeconds int64) error {
	return r.Client.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}
