package redis

import (
	"context"
	"fmt"
	"go-backend-project/internal/config"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type RedisService struct {
	RedisClient *redis.Client
}

func NewRedisConnection(cfg *config.Config) (*RedisService, error) {
	config := New(
		NewHost(cfg.RedisConfig.Host),
		NewPassword(cfg.RedisConfig.Password),
		NewPort(cfg.RedisConfig.Port),
		NewDB(cfg.RedisConfig.DB),
	)
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     30,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	ctx := context.Background()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Could not connect to Redis:", err)
		return nil, err
	}
	fmt.Println("Connected to Redis:", pong)

	redisService := RedisService{
		RedisClient: rdb,
	}

	return &redisService, nil
}
