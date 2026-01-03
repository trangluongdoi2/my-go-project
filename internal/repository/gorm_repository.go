package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type GormRepository[T any] struct {
	db *gorm.DB
}

func NewGormRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &GormRepository[T]{db: db}
}

func (r *GormRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *GormRepository[T]) CreateMany(entities []T) error {
	return r.db.Create(&entities).Error
}

func (r *GormRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

func (r *GormRepository[T]) UpdateFields(entity *T, fields map[string]interface{}) error {
	return r.db.Model(entity).Updates(fields).Error
}

func (r *GormRepository[T]) Delete(id interface{}) error {
	var entity T
	return r.db.Delete(&entity, id).Error
}

func (r *GormRepository[T]) SoftDelete(id interface{}) error {
	return r.Delete(id)
}

func (r *GormRepository[T]) HardDelete(id interface{}) error {
	var entity T
	return r.db.Unscoped().Delete(&entity, id).Error
}

func (r *GormRepository[T]) FindByID(id interface{}) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GormRepository[T]) FindOne(condition interface{}) (*T, error) {
	var entity T
	err := r.db.Where(condition).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GormRepository[T]) FindOneWithPreload(condition interface{}, preloads ...string) (*T, error) {
	var entity T
	query := r.db.Where(condition)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	err := query.First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GormRepository[T]) FindAll() ([]T, error) {
	var entities []T
	err := r.db.Find(&entities).Error
	return entities, err
}

func (r *GormRepository[T]) FindBy(condition interface{}) ([]T, error) {
	var entities []T
	err := r.db.Where(condition).Find(&entities).Error
	return entities, err
}

func (r *GormRepository[T]) FindByWithPreload(condition interface{}, preloads ...string) ([]T, error) {
	var entities []T
	query := r.db.Where(condition)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	err := query.Find(&entities).Error
	return entities, err
}

func (r *GormRepository[T]) FindWithPagination(offset, limit int, condition interface{}) ([]T, int64, error) {
	var entities []T
	var total int64

	query := r.db.Model(new(T))

	if condition != nil {
		query = query.Where(condition)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

func (r *GormRepository[T]) FindWithSort(condition interface{}, orderBy string) ([]T, error) {
	var entities []T
	query := r.db

	if condition != nil {
		query = query.Where(condition)
	}

	err := query.Order(orderBy).Find(&entities).Error
	return entities, err
}

func (r *GormRepository[T]) Search(query string, fields []string) ([]T, error) {
	var entities []T

	if len(fields) == 0 {
		return entities, fmt.Errorf("search fields cannot be empty")
	}

	db := r.db
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

func (r *GormRepository[T]) Count(condition interface{}) (int64, error) {
	var count int64
	query := r.db.Model(new(T))

	if condition != nil {
		query = query.Where(condition)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *GormRepository[T]) Exists(condition interface{}) (bool, error) {
	count, err := r.Count(condition)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GormRepository[T]) WithTransaction(tx *gorm.DB) BaseRepository[T] {
	return &GormRepository[T]{db: tx}
}

func (r *GormRepository[T]) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *GormRepository[T]) FindByDateRange(field string, startDate, endDate time.Time) ([]T, error) {
	var entities []T
	err := r.db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", field), startDate, endDate).Find(&entities).Error
	return entities, err
}

func (r *GormRepository[T]) FindByDate(field string, date time.Time) ([]T, error) {
	var entities []T
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.db.Where(fmt.Sprintf("%s >= ? AND %s < ?", field, field), startOfDay, endOfDay).Find(&entities).Error
	return entities, err
}

func (r *GormRepository[T]) FindBefore(field string, date time.Time) ([]T, error) {
	var entities []T
	err := r.db.Where(fmt.Sprintf("%s < ?", field), date).Find(&entities).Error
	return entities, err
}

func (r *GormRepository[T]) FindAfter(field string, date time.Time) ([]T, error) {
	var entities []T
	err := r.db.Where(fmt.Sprintf("%s > ?", field), date).Find(&entities).Error
	return entities, err
}
