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

// @title Go Backend Project API
// @version 1.0
// @description API Server for Go Backend Project
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:9512
// @BasePath /
// @schemes http https

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

	limiter := ratelimit.NewRedisLimiter(rdb, config.RATE_LIMITER, time.Minute)

	r := gin.Default()

	r.Use(middleware.RateLimit(limiter))

	jwtService := auth.NewJWTService(cfg.JWTConfig.Secret, cfg.JWTConfig.ExpirationHours)

	userRepo := user.NewRepository(databaseService.DB)
	authService := auth.NewService(userRepo, jwtService)
	authHandler := auth.NewHandler(authService)
	authHandler.RegisterRoutes(r)

	healthHandler := health.NewHandler(databaseService.DB, mq)
	r.GET("/health", healthHandler.HealthCheck)

	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(r)

	appointmentRepo := appointment.NewRepository(databaseService.DB)
	appointmentService := appointment.NewService(appointmentRepo, mq)
	appointmentHandler := appointment.NewHandler(appointmentService)
	appointmentHandler.RegisterRoutes(r)

	staffRepo := staff.NewRepository(databaseService.DB)
	staffService := staff.NewService(staffRepo)
	staffHandler := staff.NewHandler(staffService)
	staffHandler.RegisterRoutes(r)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtService))
	{
		serviceOfferingRepo := serviceoffering.NewRepository(databaseService.DB)
		serviceOfferingService := serviceoffering.NewService(serviceOfferingRepo)
		serviceOfferingHandler := serviceoffering.NewHandler(serviceOfferingService)
		serviceOfferingHandler.RegisterRoutes(protected)

		queueService := queue.NewService(mq)
		queueHandler := queue.NewHandler(queueService)
		queueHandler.RegisterRoutes(protected)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	r.Run(addr)
}
