package booking

import (
	"encoding/json"
	"time"

	"go-backend-project/internal/config"
	"go-backend-project/internal/rabbitmq"
)

type Service interface {
	Create(payload Booking) (*Booking, error)
	List() ([]Booking, error)
}

type service struct {
	repo      Repository
	mq        *rabbitmq.AmqpQueueService
	validator *Validator
}

func NewService(repo Repository, mq *rabbitmq.AmqpQueueService) Service {
	return &service{
		repo:      repo,
		mq:        mq,
		validator: NewValidator(),
	}
}

func (s *service) Create(payload Booking) (*Booking, error) {
	if err := s.validator.ValidateBooking(&payload); err != nil {
		return nil, err
	}

	payload.Status = int16(config.BOOKING_STATUS_PENDING)
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()

	if err := s.repo.Create(&payload); err != nil {
		return nil, err
	}

	if s.mq != nil {
		bookingData, _ := json.Marshal(payload)
		s.mq.Send("booking.created", bookingData)
	}

	return &payload, nil
}

func (s *service) List() ([]Booking, error) {
	return s.repo.GetBooking(nil)
}
