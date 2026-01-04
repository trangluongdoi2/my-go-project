package staff

import (
	"context"

	baserepo "go-backend-project/internal/repository"

	"gorm.io/gorm"
)

type Repository interface {
	baserepo.BaseRepository[Staff]

	GetStaff(ctx context.Context, condition map[string]interface{}) ([]Staff, error)
}

type repository struct {
	baserepo.BaseRepository[Staff]
	db *gorm.DB
}

func (r *repository) GetStaff(ctx context.Context, condition map[string]interface{}) ([]Staff, error) {
	return r.FindBy(ctx, condition)
}

func NewRepository(db *gorm.DB) Repository {
	baseRepo := baserepo.NewGormRepository[Staff](db)
	return &repository{
		BaseRepository: baseRepo,
		db:             db,
	}
}
