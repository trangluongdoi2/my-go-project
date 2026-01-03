package queue

import (
	"go-backend-project/internal/rabbitmq"

	"github.com/streadway/amqp"
)

type Service interface {
	PostQueue(queueName string, queuePayload []byte) (int, error)
	GetQueueOne(queueName string) (*amqp.Delivery, error)
}

type service struct {
	mq *rabbitmq.AmqpQueueService
}

func NewService(mq *rabbitmq.AmqpQueueService) Service {
	return &service{mq: mq}
}

func (s *service) PostQueue(queueName string, queuePayload []byte) (int, error) {
	err := s.mq.Send(queueName, queuePayload)
	if err != nil {
		return 500, err
	}
	return 200, nil
}

func (s *service) GetQueueOne(queueName string) (*amqp.Delivery, error) {
	msg, err := s.mq.GetOne(queueName)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
