package staff

type Service interface {
	List() ([]Staff, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) List() ([]Staff, error) {
	return s.repo.GetStaff(nil)
}
