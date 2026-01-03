package staff

import (
	baserepo "go-backend-project/internal/repository"

	"gorm.io/gorm"
)

type Repository interface {
	baserepo.BaseRepository[Staff]

	GetStaff(condition interface{}) ([]Staff, error)
}

type repository struct {
	baserepo.BaseRepository[Staff]
	db *gorm.DB
}

func (r *repository) GetStaff(condition interface{}) ([]Staff, error) {
	return r.FindBy(condition)
}
