package user

import (
	"context"
	"errors"
	"time"

	"go-backend-project/internal/config"
	"go-backend-project/pkg"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(ctx context.Context, payload User) (*User, error)
	GetUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, id string, payload User) (*User, error)
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

func (s *service) Create(ctx context.Context, payload User) (*User, error) {
	existing, _ := s.repo.GetUserByEmail(ctx, payload.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	payload.Password = string(hashedPassword)
	payload.IsActive = true
	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()
	payload.Code = pkg.GenerateCode(config.PREFIX_USER)

	if payload.Role == "" {
		payload.Role = "customer"
	}

	if err := s.repo.Create(ctx, &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func (s *service) GetUsers(ctx context.Context) ([]User, error) {
	return s.repo.GetUsers(ctx, nil)
}

func (s *service) GetUserByID(ctx context.Context, id string) (*User, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetUserByID(ctx, uuid)
}

func (s *service) Update(ctx context.Context, id string, payload User) (*User, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.GetUserByID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	payload.ID = existing.ID
	payload.Code = existing.Code
	payload.Password = existing.Password
	payload.CreatedAt = existing.CreatedAt
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
