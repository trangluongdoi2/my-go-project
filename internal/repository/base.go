package repository

import (
	"context"

	"gorm.io/gorm"
)

type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	CreateMany(ctx context.Context, entities []T) error
	Update(ctx context.Context, entity *T) error
	UpdateFields(ctx context.Context, entity *T, fields map[string]interface{}) error
	Delete(ctx context.Context, id interface{}) error
	SoftDelete(ctx context.Context, id interface{}) error
	HardDelete(ctx context.Context, id interface{}) error

	FindByID(ctx context.Context, id interface{}) (*T, error)
	FindOne(ctx context.Context, condition map[string]interface{}) (*T, error)
	FindOneWithPreload(ctx context.Context, condition map[string]interface{}, preloads []PreloadOption) (*T, error)

	FindAll(ctx context.Context) ([]T, error)
	FindBy(ctx context.Context, condition map[string]interface{}) ([]T, error)
	FindByWithPreload(ctx context.Context, condition map[string]interface{}, preloads []PreloadOption) ([]T, error)

	Find(ctx context.Context, conditions map[string]interface{}, options QueryOptions) ([]T, error)
	FindWithPreload(ctx context.Context, preloads []PreloadOption, selectFields []string, conditions map[string]interface{}, limit, offset int, sortFields []string, sortOrders []string) ([]T, error)

	Paginate(ctx context.Context, conditions map[string]interface{}, page, limit int, sortFields []string, sortOrders []string) (Pagination[T], error)
	PaginateWithPreload(ctx context.Context, conditions map[string]interface{}, page, limit int, sortFields []string, sortOrders []string, preloads []PreloadOption) (Pagination[T], error)

	FindWithSort(ctx context.Context, condition map[string]interface{}, orderBy string) ([]T, error)
	Search(ctx context.Context, query string, fields []string) ([]T, error)

	Count(ctx context.Context, condition map[string]interface{}) (int64, error)
	Exists(ctx context.Context, condition map[string]interface{}) (bool, error)

	FindByOr(ctx context.Context, conditions []map[string]interface{}) ([]T, error)
	Upsert(ctx context.Context, entity *T, conflictColumns []string, updateColumns []string) error

	WithTransaction(tx *gorm.DB) BaseRepository[T]
	BeginTransaction() *gorm.DB
}
