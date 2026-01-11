package user

import (
	"context"

	baserepo "go-backend-project/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	baserepo.BaseRepository[User]

	GetUsers(ctx context.Context, condition map[string]interface{}) ([]User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct {
	baserepo.BaseRepository[User]
	db *gorm.DB
}

func (r *repository) GetUsers(ctx context.Context, condition map[string]interface{}) ([]User, error) {
	return r.FindBy(ctx, condition)
}

func (r *repository) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return r.FindByID(ctx, id)
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return r.FindOne(ctx, map[string]interface{}{"email": email})
}

func NewRepository(db *gorm.DB) Repository {
	baseRepo := baserepo.NewGormRepository[User](db)
	return &repository{
		BaseRepository: baseRepo,
		db:             db,
	}
}
