package config

type BookingStatus int16

const (
	RATE_LIMITER = 100

	BOOKING_STATUS_PENDING   BookingStatus = 1
	BOOKING_STATUS_CONFIRMED BookingStatus = 2
	BOOKING_STATUS_COMPLETED BookingStatus = 3
	BOOKING_STATUS_CANCELLED BookingStatus = 4
)
