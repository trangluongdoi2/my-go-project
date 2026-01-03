package repository

import (
	"gorm.io/gorm"
)

type BaseRepository[T any] interface {
	Create(entity *T) error
	CreateMany(entities []T) error
	Update(entity *T) error
	UpdateFields(entity *T, fields map[string]interface{}) error
	Delete(id interface{}) error
	SoftDelete(id interface{}) error
	HardDelete(id interface{}) error

	FindByID(id interface{}) (*T, error)
	FindOne(condition interface{}) (*T, error)
	FindOneWithPreload(condition interface{}, preloads ...string) (*T, error)

	FindAll() ([]T, error)
	FindBy(condition interface{}) ([]T, error)
	FindByWithPreload(condition interface{}, preloads ...string) ([]T, error)

	FindWithPagination(offset, limit int, condition interface{}) ([]T, int64, error)
	FindWithSort(condition interface{}, orderBy string) ([]T, error)
	Search(query string, fields []string) ([]T, error)

	Count(condition interface{}) (int64, error)
	Exists(condition interface{}) (bool, error)

	WithTransaction(tx *gorm.DB) BaseRepository[T]
	BeginTransaction() *gorm.DB
}
