package booking

import (
	"context"
	"encoding/json"
	"time"

	"go-backend-project/internal/config"
	"go-backend-project/internal/rabbitmq"
	"go-backend-project/pkg"
)

type Service interface {
	Create(ctx context.Context, payload Booking) (*Booking, error)
	GetBookings(ctx context.Context) ([]Booking, error)
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

func (s *service) Create(ctx context.Context, payload Booking) (*Booking, error) {
	if err := s.validator.ValidateBooking(&payload); err != nil {
		return nil, err
	}

	payload.Status = int16(config.BOOKING_STATUS_PENDING)
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()
	payload.Code = pkg.GenerateCode(config.PREFIX_BOOKING)

	if err := s.repo.Create(ctx, &payload); err != nil {
		return nil, err
	}

	if s.mq != nil {
		bookingData, _ := json.Marshal(payload)
		s.mq.Send("booking.created", bookingData)
	}

	return &payload, nil
}

func (s *service) GetBookings(ctx context.Context) ([]Booking, error) {
	return s.repo.GetBookings(ctx, nil)
}
