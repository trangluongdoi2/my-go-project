package main

import (
	"fmt"
	"go-backend-project/internal/appointment"
	"go-backend-project/internal/auth"
	"go-backend-project/internal/config"
	"go-backend-project/internal/db"
	"go-backend-project/internal/health"
	"go-backend-project/internal/middleware"
	"go-backend-project/internal/queue"
	"go-backend-project/internal/rabbitmq"
	"go-backend-project/internal/ratelimit"
	"go-backend-project/internal/redis"
	serviceoffering "go-backend-project/internal/service-offering"
	"go-backend-project/internal/staff"
	"go-backend-project/internal/user"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go-backend-project/docs"
)

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

	limiter := ratelimit.NewRedisLimiter(rdb.RedisClient, config.RATE_LIMITER, time.Minute)

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimit(limiter))

	jwtService := auth.NewJWTService(
		cfg.JWTConfig.Secret,
		cfg.JWTConfig.AccessTokenExpirationMin,
		cfg.JWTConfig.RefreshTokenExpirationDay,
	)

	userRepo := user.NewRepository(databaseService.DB)
	authService := auth.NewService(userRepo, jwtService, rdb)
	authHandler := auth.NewHandler(authService)
	authHandler.RegisterPublicRoutes(r)

	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(r)

	appointmentRepo := appointment.NewRepository(databaseService.DB)
	appointmentService := appointment.NewService(appointmentRepo, mq)
	appointmentHandler := appointment.NewHandler(appointmentService)
	appointmentHandler.RegisterRoutes(r)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtService))
	{
		authHandler.RegisterProtectedRoutes(protected)

		staffRepo := staff.NewRepository(databaseService.DB)
		staffService := staff.NewService(staffRepo)
		staffHandler := staff.NewHandler(staffService)
		staffHandler.RegisterRoutes(r)

		serviceOfferingRepo := serviceoffering.NewRepository(databaseService.DB)
		serviceOfferingService := serviceoffering.NewService(serviceOfferingRepo)
		serviceOfferingHandler := serviceoffering.NewHandler(serviceOfferingService)
		serviceOfferingHandler.RegisterRoutes(protected)

		queueService := queue.NewService(mq)
		queueHandler := queue.NewHandler(queueService)
		queueHandler.RegisterRoutes(protected)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	healthHandler := health.NewHandler(databaseService.DB, mq)
	r.GET("/health", healthHandler.HealthCheck)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	r.Run(addr)
}
