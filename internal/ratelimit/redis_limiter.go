package ratelimit

import (
	"context"
	"fmt"
	numberutil "go-backend-project/pkg/number-util"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func NewRedisLimiter(client *redis.Client, limit int, window time.Duration) *RedisLimiter {
	return &RedisLimiter{
		client: client,
		limit:  limit,
		window: window,
	}
}

var LUA_SCRIPT = redis.NewScript(`
	local current = redis.call("INCR", KEYS[1])
	local ttl = redis.call("TTL", KEYS[1])

	if current == 1 then
		redis.call("EXPIRE", KEYS[1], ARGV[1])
		ttl = ARGV[1]
	end

	local remaining = tonumber(ARGV[2]) - current
	if remaining < 0 then
		remaining = 0
	end

	return {
		current,
		remaining,
		ttl
	}
`)

func (r *RedisLimiter) Allow(ctx context.Context, key string) (*ResultRateLimit, error) {
	redisKey := fmt.Sprintf("rate_limit:%s", key)

	res, err := LUA_SCRIPT.Run(
		ctx,
		r.client,
		[]string{redisKey},      // KEY[]
		int(r.window.Seconds()), // ARGS[1]
		r.limit,                 // ARGS[2]
	).Slice()

	if err != nil {
		return nil, err
	}
	current, err := numberutil.ToInt64(res[0])
	if err != nil {
		return nil, err
	}

	remaining, err := numberutil.ToInt64(res[1])
	if err != nil {
		return nil, err
	}

	ttl, err := numberutil.ToInt64(res[2])
	if err != nil {
		return nil, err
	}
	return &ResultRateLimit{
		Allowed:   int(current) <= r.limit,
		Limit:     r.limit,
		Remaining: int(remaining),
		ResetAt:   time.Now().Add(time.Duration(ttl) * time.Second).Unix(),
	}, nil

}
