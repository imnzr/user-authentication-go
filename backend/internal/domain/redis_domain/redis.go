package redisdomain

import "context"

type Client interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, key string, value string, ttlSecond int64) error
	Get(ctx context.Context, key string) (string, error)
}
