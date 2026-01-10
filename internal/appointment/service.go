package appointment

import (
	"context"
	"encoding/json"
	"time"

	"go-backend-project/internal/config"
	"go-backend-project/internal/rabbitmq"
	"go-backend-project/pkg"
)

type Service interface {
	Create(ctx context.Context, payload Appointment) (*Appointment, error)
	GetAppointments(ctx context.Context) ([]Appointment, error)
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

func (s *service) Create(ctx context.Context, payload Appointment) (*Appointment, error) {
	if err := s.validator.ValidateAppointment(&payload); err != nil {
		return nil, err
	}

	payload.Status = int16(config.APPOINTMENT_STATUS_PENDING)
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()
	payload.Code = pkg.GenerateCode(config.PREFIX_APPOINTMENT)

	if err := s.repo.Create(ctx, &payload); err != nil {
		return nil, err
	}

	if s.mq != nil {
		AppointmentData, _ := json.Marshal(payload)
		s.mq.Send("Appointment.created", AppointmentData)
	}

	return &payload, nil
}

func (s *service) GetAppointments(ctx context.Context) ([]Appointment, error) {
	return s.repo.GetAppointments(ctx, nil)
}
