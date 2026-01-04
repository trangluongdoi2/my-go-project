package booking

import (
	"context"
	"time"

	baserepo "go-backend-project/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	baserepo.BaseRepository[Booking]

	GetBookings(ctx context.Context, condition map[string]interface{}) ([]Booking, error)
	GetBookingByID(ctx context.Context, id uuid.UUID) (*Booking, error)
	GetBookingsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]Booking, error)
}

type repository struct {
	baserepo.BaseRepository[Booking]
	db *gorm.DB
}

func (r *repository) GetBookings(ctx context.Context, condition map[string]interface{}) ([]Booking, error) {
	return r.FindBy(ctx, condition)
}

func (r *repository) GetBookingByID(ctx context.Context, id uuid.UUID) (*Booking, error) {
	return r.FindByID(ctx, id)
}

func (r *repository) GetBookingsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]Booking, error) {
	var bookings []Booking
	err := r.db.WithContext(ctx).Where("appointment_at BETWEEN ? AND ?", startDate, endDate).
		Order("appointment_at ASC").
		Find(&bookings).Error
	return bookings, err
}

func NewRepository(db *gorm.DB) Repository {
	baseRepo := baserepo.NewGormRepository[Booking](db)
	return &repository{
		BaseRepository: baseRepo,
		db:             db,
	}
}
