package booking

import "github.com/google/uuid"

type Booking struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
