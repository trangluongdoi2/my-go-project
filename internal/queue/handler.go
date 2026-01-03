package queue

import (
	"encoding/json"

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

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	queue := r.Group("/queue")
	{
		queue.GET("/execute", h.GetQueueOne)
		queue.POST("/send", h.PostQueue)
	}
}

func (h *Handler) PostQueue(c *gin.Context) {
	queueName := c.Query("queue")
	var req QueueMessage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	out, err := json.Marshal(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	status, err := h.service.PostQueue(queueName, out)

	if err != nil {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.String(status, "Send queue successfully!")
}

func (h *Handler) GetQueueOne(c *gin.Context) {
	queueName := c.Query("queue")
	if queueName == "" {
		c.JSON(400, gin.H{"error": "queue parameter is required"})
		return
	}

	msg, err := h.service.GetQueueOne(queueName)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if msg == nil {
		c.JSON(200, gin.H{"message": "Queue is empty"})
		return
	}

	if err := msg.Ack(false); err != nil {
		c.JSON(500, gin.H{"error": "Failed to acknowledge message: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message_id": msg.MessageId,
		"body":       string(msg.Body),
		"timestamp":  msg.Timestamp,
	})
}
