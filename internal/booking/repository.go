package booking

import (
	"time"

	baserepo "go-backend-project/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	baserepo.BaseRepository[Booking]

	GetBooking(condition interface{}) ([]Booking, error)
	GetBookingByID(id uuid.UUID) (*Booking, error)
	GetBookingsByStatus(status int16) ([]Booking, error)
	GetBookingsByCustomerPhone(phone string) ([]Booking, error)
	GetBookingsByStaffID(staffID uuid.UUID) ([]Booking, error)
	GetBookingsByDateRange(startDate, endDate time.Time) ([]Booking, error)
}

type repository struct {
	baserepo.BaseRepository[Booking]
	db *gorm.DB
}

func (r *repository) GetBooking(condition interface{}) ([]Booking, error) {
	return r.FindBy(condition)
}

func (r *repository) GetBookingByID(id uuid.UUID) (*Booking, error) {
	return r.FindByID(id)
}

func (r *repository) GetBookingsByStatus(status int16) ([]Booking, error) {
	var bookings []Booking
	err := r.db.Where("status = ?", status).Find(&bookings).Error
	return bookings, err
}

func (r *repository) GetBookingsByCustomerPhone(phone string) ([]Booking, error) {
	var bookings []Booking
	err := r.db.Where("customer_phone = ?", phone).Order("appointment_at DESC").Find(&bookings).Error
	return bookings, err
}

func (r *repository) GetBookingsByStaffID(staffID uuid.UUID) ([]Booking, error) {
	var bookings []Booking
	err := r.db.Where("staff_id = ?", staffID).Order("appointment_at DESC").Find(&bookings).Error
	return bookings, err
}

func (r *repository) GetBookingsByDateRange(startDate, endDate time.Time) ([]Booking, error) {
	var bookings []Booking
	err := r.db.Where("appointment_at BETWEEN ? AND ?", startDate, endDate).
		Order("appointment_at ASC").
		Find(&bookings).Error
	return bookings, err
}

// func (r *repository) GetUpcomingBookings(limit int) ([]Booking, error) {
// 	var bookings []Booking
// 	err := r.db.Where("appointment_at > ? AND status IN ?", time.Now(), []BookingStatus{BookingStatusPending, BookingStatusConfirmed}).
// 		Order("appointment_at ASC").
// 		Limit(limit).
// 		Find(&bookings).Error
// 	return bookings, err
// }

// func (r *repository) UpdateStatus(id uuid.UUID, status BookingStatus) error {
// 	return r.db.Model(&Booking{}).Where("id = ?", id).Update("status", status).Error
// }

// func (r *repository) CheckConflict(staffID uuid.UUID, appointmentAt time.Time, duration int, excludeID *uuid.UUID) (bool, error) {
// 	endTime := appointmentAt.Add(time.Duration(duration) * time.Minute)

// 	query := r.db.Model(&Booking{}).
// 		Where("staff_id = ?", staffID).
// 		Where("status IN ?", []BookingStatus{BookingStatusPending, BookingStatusConfirmed}).
// 		Where("appointment_at < ?", endTime).
// 		Where("DATE_ADD(appointment_at, INTERVAL duration MINUTE) > ?", appointmentAt)

// 	if excludeID != nil {
// 		query = query.Where("id != ?", *excludeID)
// 	}

// 	var count int64
// 	err := query.Count(&count).Error
// 	if err != nil {
// 		return false, err
// 	}

// 	return count > 0, nil
// }

func NewRepository(db *gorm.DB) Repository {
	baseRepo := baserepo.NewGormRepository[Booking](db)
	return &repository{
		BaseRepository: baseRepo,
		db:             db,
	}
}
