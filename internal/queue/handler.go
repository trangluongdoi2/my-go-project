package queue

import (
	"encoding/json"
	"go-backend-project/utils"

	"github.com/gin-gonic/gin"
)

type QueueMessage struct {
	Message string
	Data    map[string]any
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(r gin.IRouter) {
	r.GET("/queue/execute", h.GetQueueOne)
	r.POST("/queue/send", h.PostQueue)
}

func (h *Handler) PostQueue(c *gin.Context) {
	queueName := c.Query("queue")
	var req QueueMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	out, err := json.Marshal(req)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	status, err := h.service.PostQueue(queueName, out)

	if err != nil {
		utils.Error(c, status, err.Error())
		return
	}
	utils.Success(c, status, "Queue sent successfully", nil)
}

func (h *Handler) GetQueueOne(c *gin.Context) {
	queueName := c.Query("queue")
	if queueName == "" {
		utils.BadRequest(c, "queue parameter is required")
		return
	}

	msg, err := h.service.GetQueueOne(queueName)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	if msg == nil {
		utils.OK(c, "Queue is empty", nil)
		return
	}

	if err := msg.Ack(false); err != nil {
		utils.InternalServerError(c, "Failed to acknowledge message: "+err.Error())
		return
	}

	utils.OK(c, "Queue message retrieved", gin.H{
		"message_id": msg.MessageId,
		"body":       string(msg.Body),
		"timestamp":  msg.Timestamp,
	})
}
