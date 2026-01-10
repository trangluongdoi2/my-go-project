package config

type AppointmentStatus int16

const (
	RATE_LIMITER = 100

	PREFIX_APPOINTMENT = "APT"
	PREFIX_STAFF       = "SAF"
	PREFIX_SERVICE     = "SEV"

	APPOINTMENT_STATUS_PENDING   AppointmentStatus = 1
	APPOINTMENT_STATUS_CONFIRMED AppointmentStatus = 2
	APPOINTMENT_STATUS_COMPLETED AppointmentStatus = 3
	APPOINTMENT_STATUS_CANCELLED AppointmentStatus = 4
)
