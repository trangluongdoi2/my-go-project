package staff

import (
	"context"
	"errors"
	"go-backend-project/internal/config"
	baserepo "go-backend-project/internal/repository"
	"go-backend-project/pkg"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	List(ctx context.Context, query ListStaffQuery) (baserepo.Pagination[Staff], error)
	GetByID(ctx context.Context, id string) (*Staff, error)
	Create(ctx context.Context, payload Staff) (*Staff, error)
	Update(ctx context.Context, id string, payload Staff) (*Staff, error)
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

func (s *service) List(ctx context.Context, query ListStaffQuery) (baserepo.Pagination[Staff], error) {
	conditions := make(map[string]any)
	if query.Code != "" {
		conditions["code"] = query.Code
	}
	if query.FirstName != "" {
		conditions["first_name"] = query.FirstName
	}
	if query.LastName != "" {
		conditions["last_name"] = query.LastName
	}
	if query.Email != "" {
		conditions["email"] = query.Email
	}
	if query.Phone != "" {
		conditions["phone"] = query.Phone
	}
	if query.IsActive != nil {
		conditions["is_active"] = *query.IsActive
	}

	return s.repo.GetStaffPaginated(
		ctx,
		conditions,
		query.GetPage(),
		query.GetLimit(),
		query.GetSortFields(),
		query.GetSortOrders(),
	)
}

func (s *service) GetByID(ctx context.Context, id string) (*Staff, error) {
	staffID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid staff ID")
	}

	staff, err := s.repo.FindByID(ctx, staffID)
	if err != nil {
		return nil, errors.New("staff not found")
	}

	return staff, nil
}

func (s *service) Create(ctx context.Context, payload Staff) (*Staff, error) {
	existing, _ := s.repo.FindOne(ctx, map[string]interface{}{"phone": payload.Phone})

	if existing != nil {
		return nil, errors.New("staff already exists")
	}

	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()
	payload.Code = pkg.GenerateCode(config.PREFIX_STAFF)

	if err := s.repo.Create(ctx, &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func (s *service) Update(ctx context.Context, id string, payload Staff) (*Staff, error) {
	staffID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid staff ID")
	}

	existing, err := s.repo.FindByID(ctx, staffID)
	if err != nil {
		return nil, errors.New("staff not found")
	}

	if payload.FirstName != "" {
		existing.FirstName = payload.FirstName
	}
	if payload.LastName != "" {
		existing.LastName = payload.LastName
	}
	if payload.Email != "" {
		existing.Email = payload.Email
	}
	if payload.Phone != "" {
		existing.Phone = payload.Phone
	}

	existing.IsActive = payload.IsActive
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	staffID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid staff ID")
	}

	_, err = s.repo.FindByID(ctx, staffID)
	if err != nil {
		return errors.New("staff not found")
	}

	return s.repo.Delete(ctx, staffID)
}
