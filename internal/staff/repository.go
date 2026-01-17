package staff

import (
	"context"

	baserepo "go-backend-project/internal/repository"

	"gorm.io/gorm"
)

type Repository interface {
	baserepo.BaseRepository[Staff]

	GetStaff(ctx context.Context, condition map[string]interface{}) ([]Staff, error)
	GetStaffPaginated(ctx context.Context, condition map[string]interface{}, page, limit int, sortFields, sortOrders []string) (baserepo.Pagination[Staff], error)
	GetStaffBy(ctx context.Context, condition map[string]interface{}) (*Staff, error)
}

type repository struct {
	baserepo.BaseRepository[Staff]
	db *gorm.DB
}

func (r *repository) GetStaff(ctx context.Context, condition map[string]interface{}) ([]Staff, error) {
	return r.FindBy(ctx, condition)
}

func (r *repository) GetStaffPaginated(ctx context.Context, condition map[string]interface{}, page, limit int, sortFields, sortOrders []string) (baserepo.Pagination[Staff], error) {
	return r.Paginate(ctx, condition, page, limit, sortFields, sortOrders)
}

func (r *repository) GetStaffBy(ctx context.Context, condition map[string]interface{}) (*Staff, error) {
	return r.FindOne(ctx, condition)
}

func NewRepository(db *gorm.DB) Repository {
	baseRepo := baserepo.NewGormRepository[Staff](db)
	return &repository{
		BaseRepository: baseRepo,
		db:             db,
	}
}
