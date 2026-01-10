package appointment

import (
	"context"
	"time"

	baserepo "go-backend-project/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	baserepo.BaseRepository[Appointment]

	GetAppointments(ctx context.Context, condition map[string]interface{}) ([]Appointment, error)
	GetAppointmentByID(ctx context.Context, id uuid.UUID) (*Appointment, error)
	GetAppointmentsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]Appointment, error)
}

type repository struct {
	baserepo.BaseRepository[Appointment]
	db *gorm.DB
}

func (r *repository) GetAppointments(ctx context.Context, condition map[string]interface{}) ([]Appointment, error) {
	return r.FindBy(ctx, condition)
}

func (r *repository) GetAppointmentByID(ctx context.Context, id uuid.UUID) (*Appointment, error) {
	return r.FindByID(ctx, id)
}

func (r *repository) GetAppointmentsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]Appointment, error) {
	var Appointments []Appointment
	err := r.db.WithContext(ctx).Where("appointment_at BETWEEN ? AND ?", startDate, endDate).
		Order("appointment_at ASC").
		Find(&Appointments).Error
	return Appointments, err
}

func NewRepository(db *gorm.DB) Repository {
	baseRepo := baserepo.NewGormRepository[Appointment](db)
	return &repository{
		BaseRepository: baseRepo,
		db:             db,
	}
}
