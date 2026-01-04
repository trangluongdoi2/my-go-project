package serviceoffering

import (
	"context"

	baserepo "go-backend-project/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	baserepo.BaseRepository[ServiceOffering]

	GetServiceOfferings(ctx context.Context, condition map[string]interface{}) ([]ServiceOffering, error)
	GetServiceOfferingByID(ctx context.Context, id uuid.UUID) (*ServiceOffering, error)

	AssignStaffToService(ctx context.Context, serviceID uuid.UUID, staffIDs []uuid.UUID) error
	RemoveStaffFromService(ctx context.Context, serviceID uuid.UUID, staffIDs []uuid.UUID) error
	GetServiceWithStaffs(ctx context.Context, serviceID uuid.UUID) (*ServiceOffering, error)
}

type repository struct {
	baserepo.BaseRepository[ServiceOffering]
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	baseRepo := baserepo.NewGormRepository[ServiceOffering](db)
	return &repository{
		BaseRepository: baseRepo,
		db:             db,
	}
}

func (r *repository) GetServiceOfferings(ctx context.Context, condition map[string]interface{}) ([]ServiceOffering, error) {
	if condition == nil {
		return r.FindAll(ctx)
	}
	return r.FindBy(ctx, condition)
}

func (r *repository) GetServiceOfferingByID(ctx context.Context, id uuid.UUID) (*ServiceOffering, error) {
	return r.FindByID(ctx, id)
}

func (r *repository) AssignStaffToService(ctx context.Context, serviceID uuid.UUID, staffIDs []uuid.UUID) error {
	service := &ServiceOffering{ID: serviceID}

	for _, staffID := range staffIDs {
		staff := Staff{ID: staffID}
		if err := r.db.WithContext(ctx).Model(service).Association("Staffs").Append(&staff); err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) RemoveStaffFromService(ctx context.Context, serviceID uuid.UUID, staffIDs []uuid.UUID) error {
	service := &ServiceOffering{ID: serviceID}

	for _, staffID := range staffIDs {
		staff := Staff{ID: staffID}
		if err := r.db.WithContext(ctx).Model(service).Association("Staffs").Delete(&staff); err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) GetServiceWithStaffs(ctx context.Context, serviceID uuid.UUID) (*ServiceOffering, error) {
	var service ServiceOffering
	err := r.db.WithContext(ctx).Preload("Staffs").First(&service, serviceID).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}
