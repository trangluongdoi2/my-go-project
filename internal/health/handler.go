package health

import (
	"go-backend-project/internal/rabbitmq"
	"go-backend-project/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
	mq *rabbitmq.AmqpQueueService
}

func NewHandler(db *gorm.DB, mq *rabbitmq.AmqpQueueService) *Handler {
	return &Handler{db: db, mq: mq}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	status := map[string]string{
		"database": "ok",
		"rabbitmq": "ok",
	}

	db, err := h.db.DB()
	if err != nil || db.Ping() != nil {
		status["database"] = "down"
	}

	utils.OK(c, "Health check completed", status)
}
