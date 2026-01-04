package staff

import "context"

type Service interface {
	List(ctx context.Context) ([]Staff, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) List(ctx context.Context) ([]Staff, error) {
	return s.repo.GetStaff(ctx, nil)
}
