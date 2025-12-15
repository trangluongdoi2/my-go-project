package main

import (
	"fmt"
	"go-backend-project/internal/booking"
	"go-backend-project/internal/config"
	"go-backend-project/internal/db"
	"go-backend-project/internal/health"
	"go-backend-project/internal/middleware"
	"go-backend-project/internal/queue"
	"go-backend-project/internal/rabbitmq"
	"go-backend-project/internal/ratelimit"
	"go-backend-project/internal/redis"
	"time"

	"github.com/gin-gonic/gin"
)

// func AA() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		fmt.Println("Before ==> AA")
// 		c.Next()
// 		fmt.Println("After ==> AA")
// 	}
// }

// func BB() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		fmt.Println("Before ==> BB")
// 		c.Next()
// 		fmt.Println("After ==> BB")
// 	}
// }
// func CC() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		fmt.Println("Before ==> CC")
// 		c.Next()
// 		fmt.Println("After ==> CC")
// 	}
// }

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err.Error(), "error when load config")
	}
	databaseService, err := db.NewPostgresConnection(cfg)
	if err != nil {
		fmt.Println(err.Error(), "error when init database")
	}

	mq, _ := rabbitmq.NewAmqpQueueService(cfg)

	rdb, _ := redis.NewRedisConnection(cfg)

	// limiter := ratelimit.NewRedisLimiter(rdb, 100, time.Minute)
	limiter := ratelimit.NewRedisLimiter(rdb, 5, time.Minute)

	r := gin.Default()

	r.Use(middleware.RateLimit(limiter))

	bookingRepo := booking.NewRepository(databaseService.DB)
	bookingService := booking.NewService(bookingRepo, mq)
	bookingHandler := booking.NewHandler(bookingService)
	bookingHandler.RegisterRoutes(r)

	healthHandler := health.NewHandler(databaseService.DB, mq)

	queueService := queue.NewService(mq)
	queueHandler := queue.NewHandler(queueService)
	queueHandler.RegisterRoutes(r)

	r.GET("/health", healthHandler.HealthCheck)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	r.Run(addr)
}
