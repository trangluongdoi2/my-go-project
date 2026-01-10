package appointment

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateAppointment(pppointment *Appointment) error {
	// if Appointment.ServiceName == "" {
	// 	return errors.New("service name is required")
	// }
	// if Appointment.AppointmentAt.IsZero() {
	// 	return errors.New("appointment date and time is required")
	// }
	// if Appointment.AppointmentAt.Before(time.Now()) {
	// 	return errors.New("appointment time must be in the future")
	// }
	// if Appointment.Price < 0 {
	// 	return errors.New("price cannot be negative")
	// }
	return nil
}
