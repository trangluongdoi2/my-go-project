package queue

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

type QueueMessage struct {
	Message string
	Data    map[string]interface{}
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
	return
}

func (h *Handler) GetQueueOne(c *gin.Context) {
	queueName := c.Query("queue")
	fmt.Println(queueName, "queueName...")
	msg, err := h.service.GetQueueOne(queueName)
	if err != nil {
		fmt.Println("error getQueueOne...")
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if msg == nil {
		c.JSON(200, gin.H{"message": "Queue empty"})
		return
	}

	msg.Ack(false)

	c.JSON(200, gin.H{
		"message_id": msg.MessageId,
		"body":       string(msg.Body),
		"timestamp":  msg.Timestamp,
	})
}
