package booking

import "gorm.io/gorm"

type Repository interface {
	Create(booking *Booking) error
	GetAll() ([]Booking, error)
}

type repository struct {
	db *gorm.DB
}

func (r *repository) Create(booking *Booking) error {
	return r.db.Create(booking).Error
}

func (r *repository) GetAll() ([]Booking, error) {
	var bookings []Booking
	err := r.db.Find(&bookings).Error
	return bookings, err
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
