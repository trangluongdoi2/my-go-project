package serviceoffering

import (
	"context"

	"go-backend-project/internal/config"
	"go-backend-project/pkg"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, payload ServiceOffering) (*ServiceOffering, error)
	GetServiceOfferings(ctx context.Context, condition map[string]interface{}) ([]ServiceOffering, error)
	GetServiceOfferingByID(ctx context.Context, id string) (*ServiceOffering, error)
	Update(ctx context.Context, id string, payload ServiceOffering) (*ServiceOffering, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, payload ServiceOffering) (*ServiceOffering, error) {
	payload.Status = 1
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()
	payload.Code = pkg.GenerateCode(config.PREFIX_SERVICE)

	if err := s.repo.Create(ctx, &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func (s *service) GetServiceOfferings(ctx context.Context, condition map[string]interface{}) ([]ServiceOffering, error) {
	return s.repo.GetServiceOfferings(ctx, condition)
}

func (s *service) GetServiceOfferingByID(ctx context.Context, id string) (*ServiceOffering, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetServiceOfferingByID(ctx, uuid)
}

func (s *service) Update(ctx context.Context, id string, payload ServiceOffering) (*ServiceOffering, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.GetServiceOfferingByID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	payload.ID = existing.ID
	payload.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.SoftDelete(ctx, uuid)
}
