package booking

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateBooking(booking *Booking) error {
	// if booking.ServiceName == "" {
	// 	return errors.New("service name is required")
	// }
	// if booking.AppointmentAt.IsZero() {
	// 	return errors.New("appointment date and time is required")
	// }
	// if booking.AppointmentAt.Before(time.Now()) {
	// 	return errors.New("appointment time must be in the future")
	// }
	// if booking.Price < 0 {
	// 	return errors.New("price cannot be negative")
	// }
	return nil
}
