package repository

import (
	"context"
	"fmt"
	"math"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormRepository[T any] struct {
	db *gorm.DB
}

func NewGormRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &GormRepository[T]{db: db}
}

func (r *GormRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *GormRepository[T]) CreateMany(ctx context.Context, entities []T) error {
	return r.db.WithContext(ctx).Create(&entities).Error
}

func (r *GormRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *GormRepository[T]) UpdateFields(ctx context.Context, entity *T, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(entity).Updates(fields).Error
}

func (r *GormRepository[T]) Delete(ctx context.Context, id interface{}) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, id).Error
}

func (r *GormRepository[T]) SoftDelete(ctx context.Context, id interface{}) error {
	return r.Delete(ctx, id)
}

func (r *GormRepository[T]) HardDelete(ctx context.Context, id interface{}) error {
	var entity T
	return r.db.WithContext(ctx).Unscoped().Delete(&entity, id).Error
}

func (r *GormRepository[T]) FindByID(ctx context.Context, id interface{}) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GormRepository[T]) FindOne(ctx context.Context, condition map[string]interface{}) (*T, error) {
	var entity T
	query := r.buildQuery(r.db.WithContext(ctx), condition)
	err := query.First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GormRepository[T]) FindOneWithPreload(ctx context.Context, condition map[string]interface{}, preloads []PreloadOption) (*T, error) {
	var entity T
	query := r.buildQuery(r.db.WithContext(ctx), condition)
	query = r.applyPreloads(query, preloads)

	err := query.First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GormRepository[T]) FindAll(ctx context.Context) ([]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).Find(&entities).Error
	return entities, err
}

func (r *GormRepository[T]) FindBy(ctx context.Context, condition map[string]interface{}) ([]T, error) {
	var entities []T
	query := r.buildQuery(r.db.WithContext(ctx), condition)
	err := query.Find(&entities).Error
	return entities, err
}

func (r *GormRepository[T]) FindByWithPreload(ctx context.Context, condition map[string]interface{}, preloads []PreloadOption) ([]T, error) {
	var entities []T
	query := r.buildQuery(r.db.WithContext(ctx), condition)
	query = r.applyPreloads(query, preloads)

	err := query.Find(&entities).Error
	return entities, err
}

func (r *GormRepository[T]) Find(ctx context.Context, conditions map[string]interface{}, options QueryOptions) ([]T, error) {
	var entities []T
	query := r.buildQuery(r.db.WithContext(ctx), conditions)

	if len(options.SelectFields) > 0 {
		query = query.Select(options.SelectFields)
	}

	if len(options.SortFields) > 0 && len(options.SortOrders) > 0 {
		for i, field := range options.SortFields {
			order := "ASC"
			if i < len(options.SortOrders) {
				order = options.SortOrders[i]
			}
			query = query.Order(fmt.Sprintf("%s %s", field, order))
		}
	}

	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}

	if options.Offset > 0 {
		query = query.Offset(options.Offset)
	}

	if len(options.Preloads) > 0 {
		query = r.applyPreloads(query, options.Preloads)
	}

	err := query.Find(&entities).Error
	return entities, err
}

func (r *GormRepository[T]) FindWithPreload(ctx context.Context, preloads []PreloadOption, selectFields []string, conditions map[string]interface{}, limit, offset int, sortFields []string, sortOrders []string) ([]T, error) {
	var entities []T
	query := r.buildQuery(r.db.WithContext(ctx), conditions)

	if len(selectFields) > 0 {
		query = query.Select(selectFields)
	}

	query = r.applyPreloads(query, preloads)

	if len(sortFields) > 0 && len(sortOrders) > 0 {
		for i, field := range sortFields {
			order := "ASC"
			if i < len(sortOrders) {
				order = sortOrders[i]
			}
			query = query.Order(fmt.Sprintf("%s %s", field, order))
		}
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&entities).Error
	return entities, err
}

func (r *GormRepository[T]) Paginate(ctx context.Context, conditions map[string]interface{}, page, limit int, sortFields []string, sortOrders []string) (Pagination[T], error) {
	var entities []T
	var total int64

	query := r.buildQuery(r.db.WithContext(ctx).Model(new(T)), conditions)

	if err := query.Count(&total).Error; err != nil {
		return Pagination[T]{}, err
	}

	offset := (page - 1) * limit

	if len(sortFields) > 0 && len(sortOrders) > 0 {
		for i, field := range sortFields {
			order := "ASC"
			if i < len(sortOrders) {
				order = sortOrders[i]
			}
			query = query.Order(fmt.Sprintf("%s %s", field, order))
		}
	}

	if err := query.Limit(limit).Offset(offset).Find(&entities).Error; err != nil {
		return Pagination[T]{}, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return Pagination[T]{
		Total:     total,
		Items:     entities,
		Limit:     limit,
		Page:      page,
		TotalPage: totalPage,
	}, nil
}

func (r *GormRepository[T]) PaginateWithPreload(ctx context.Context, conditions map[string]interface{}, page, limit int, sortFields []string, sortOrders []string, preloads []PreloadOption) (Pagination[T], error) {
	var entities []T
	var total int64

	query := r.buildQuery(r.db.WithContext(ctx).Model(new(T)), conditions)

	if err := query.Count(&total).Error; err != nil {
		return Pagination[T]{}, err
	}

	offset := (page - 1) * limit

	query = r.applyPreloads(query, preloads)

	if len(sortFields) > 0 && len(sortOrders) > 0 {
		for i, field := range sortFields {
			order := "ASC"
			if i < len(sortOrders) {
				order = sortOrders[i]
			}
			query = query.Order(fmt.Sprintf("%s %s", field, order))
		}
	}

	if err := query.Limit(limit).Offset(offset).Find(&entities).Error; err != nil {
		return Pagination[T]{}, err
	}

	totalPage := int(math.Ceil(float64(total) / float64(limit)))

	return Pagination[T]{
		Total:     total,
		Items:     entities,
		Limit:     limit,
		Page:      page,
		TotalPage: totalPage,
	}, nil
}

func (r *GormRepository[T]) FindWithSort(ctx context.Context, condition map[string]interface{}, orderBy string) ([]T, error) {
	var entities []T
	query := r.buildQuery(r.db.WithContext(ctx), condition)
	err := query.Order(orderBy).Find(&entities).Error
	return entities, err
}

func (r *GormRepository[T]) Search(ctx context.Context, query string, fields []string) ([]T, error) {
	var entities []T

	if len(fields) == 0 {
		return entities, fmt.Errorf("search fields cannot be empty")
	}

	db := r.db.WithContext(ctx)
	for i, field := range fields {
		if i == 0 {
			db = db.Where(fmt.Sprintf("%s LIKE ?", field), "%"+query+"%")
		} else {
			db = db.Or(fmt.Sprintf("%s LIKE ?", field), "%"+query+"%")
		}
	}

	err := db.Find(&entities).Error
	return entities, err
}

func (r *GormRepository[T]) Count(ctx context.Context, condition map[string]interface{}) (int64, error) {
	var count int64
	query := r.buildQuery(r.db.WithContext(ctx).Model(new(T)), condition)
	err := query.Count(&count).Error
	return count, err
}

func (r *GormRepository[T]) Exists(ctx context.Context, condition map[string]interface{}) (bool, error) {
	count, err := r.Count(ctx, condition)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormRepository[T]) FindByOr(ctx context.Context, conditions []map[string]interface{}) ([]T, error) {
	var entities []T
	query := r.db.WithContext(ctx)

	for i, condition := range conditions {
		if i == 0 {
			query = r.buildQuery(query, condition)
		} else {
			query = query.Or(r.buildQuery(r.db.WithContext(ctx), condition))
		}
	}

	err := query.Find(&entities).Error
	return entities, err
}

func (r *GormRepository[T]) Upsert(ctx context.Context, entity *T, conflictColumns []string, updateColumns []string) error {
	conflictClause := clause.OnConflict{
		Columns:   make([]clause.Column, len(conflictColumns)),
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}

	for i, col := range conflictColumns {
		conflictClause.Columns[i] = clause.Column{Name: col}
	}

	return r.db.WithContext(ctx).Clauses(conflictClause).Create(entity).Error
}

func (r *GormRepository[T]) WithTransaction(tx *gorm.DB) BaseRepository[T] {
	return &GormRepository[T]{db: tx}
}

func (r *GormRepository[T]) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *GormRepository[T]) buildQuery(db *gorm.DB, conditions map[string]interface{}) *gorm.DB {
	if conditions == nil {
		return db
	}

	for key, value := range conditions {
		if strings.HasSuffix(key, "_IN") {
			field := strings.TrimSuffix(key, "_IN")
			db = db.Where(fmt.Sprintf("%s IN ?", field), value)
		} else if strings.HasSuffix(key, "_NOT_IN") {
			field := strings.TrimSuffix(key, "_NOT_IN")
			db = db.Where(fmt.Sprintf("%s NOT IN ?", field), value)
		} else if strings.HasSuffix(key, "_LIKE") {
			field := strings.TrimSuffix(key, "_LIKE")
			db = db.Where(fmt.Sprintf("%s LIKE ?", field), fmt.Sprintf("%%%v%%", value))
		} else if strings.HasSuffix(key, "_NOT_LIKE") {
			field := strings.TrimSuffix(key, "_NOT_LIKE")
			db = db.Where(fmt.Sprintf("%s NOT LIKE ?", field), fmt.Sprintf("%%%v%%", value))
		} else if strings.HasSuffix(key, "_BETWEEN_AND") {
			field := strings.TrimSuffix(key, "_BETWEEN_AND")
			if arr, ok := value.([]interface{}); ok && len(arr) == 2 {
				db = db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", field), arr[0], arr[1])
			}
		} else if strings.HasSuffix(key, "_NOT_NULL") {
			field := strings.TrimSuffix(key, "_NOT_NULL")
			db = db.Where(fmt.Sprintf("%s IS NOT NULL", field))
		} else if strings.HasSuffix(key, "_NULL") {
			field := strings.TrimSuffix(key, "_NULL")
			db = db.Where(fmt.Sprintf("%s IS NULL", field))
		} else if strings.HasSuffix(key, "_GREATER_THAN") {
			field := strings.TrimSuffix(key, "_GREATER_THAN")
			db = db.Where(fmt.Sprintf("%s > ?", field), value)
		} else if strings.HasSuffix(key, "_GREATER_THAN_OR_EQUAL") {
			field := strings.TrimSuffix(key, "_GREATER_THAN_OR_EQUAL")
			db = db.Where(fmt.Sprintf("%s >= ?", field), value)
		} else if strings.HasSuffix(key, "_LESS_THAN") {
			field := strings.TrimSuffix(key, "_LESS_THAN")
			db = db.Where(fmt.Sprintf("%s < ?", field), value)
		} else if strings.HasSuffix(key, "_LESS_THAN_OR_EQUAL") {
			field := strings.TrimSuffix(key, "_LESS_THAN_OR_EQUAL")
			db = db.Where(fmt.Sprintf("%s <= ?", field), value)
		} else if strings.HasSuffix(key, "_NOT_EQUAL") {
			field := strings.TrimSuffix(key, "_NOT_EQUAL")
			db = db.Where(fmt.Sprintf("%s != ?", field), value)
		} else if strings.HasSuffix(key, "_EQUAL") {
			field := strings.TrimSuffix(key, "_EQUAL")
			db = db.Where(fmt.Sprintf("%s = ?", field), value)
		} else {
			db = db.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	return db
}

func (r *GormRepository[T]) applyPreloads(db *gorm.DB, preloads []PreloadOption) *gorm.DB {
	for _, preload := range preloads {
		if len(preload.Conditions) > 0 || len(preload.SelectFields) > 0 || preload.Limit > 0 || preload.Offset > 0 {
			db = db.Preload(preload.Entity, func(db *gorm.DB) *gorm.DB {
				if len(preload.Conditions) > 0 {
					db = r.buildQuery(db, preload.Conditions)
				}
				if len(preload.SelectFields) > 0 {
					db = db.Select(preload.SelectFields)
				}
				if preload.Limit > 0 {
					db = db.Limit(preload.Limit)
				}
				if preload.Offset > 0 {
					db = db.Offset(preload.Offset)
				}
				if len(preload.NestedPreloads) > 0 {
					db = r.applyPreloads(db, preload.NestedPreloads)
				}
				return db
			})
		} else {
			db = db.Preload(preload.Entity)
			if len(preload.NestedPreloads) > 0 {
				db = r.applyPreloads(db, preload.NestedPreloads)
			}
		}
	}
	return db
}
