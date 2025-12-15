package booking

import "go-backend-project/internal/rabbitmq"

type Service interface {
	Create(payload Booking) (int, error)
	List() ([]Booking, error)
}

type service struct {
	repo Repository
	mq   *rabbitmq.AmqpQueueService
}

func NewService(repo Repository, mq *rabbitmq.AmqpQueueService) Service {
	return &service{repo: repo, mq: mq}
}

func (s *service) Create(payload Booking) (int, error) {
	return 0, nil
}

func (s *service) List() ([]Booking, error) {
	return s.repo.GetAll()
}
