package redis

import "context"

type Client interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, key string, value string, ttlSeconds int64) error
	Get(ctx context.Context, key string) (string, error)
}
