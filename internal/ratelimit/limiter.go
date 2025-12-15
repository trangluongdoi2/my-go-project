package ratelimit

import "context"

type Limiter interface {
	Allow(ctx context.Context, key string) (*ResultRateLimit, error)
}
