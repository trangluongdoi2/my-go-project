package appointment

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Appointment struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Code          string         `json:"code"`
	Customer      datatypes.JSON `json:"customer" gorm:"type:jsonb" swaggertype:"object"`
	ServiceName   string         `json:"service_name" gorm:"not null"`
	StaffID       *uuid.UUID     `json:"staff_id" gorm:"type:uuid"`
	AppointmentAt time.Time      `json:"appointment_at"`
	Price         float64        `json:"price" gorm:"type:decimal(10,2)"`
	Status        int16          `json:"status" gorm:"type:smallint;default:1"`
	PaymentStatus int16          `json:"payment_status" gorm:"type:smallint;default:1"`
	Notes         string         `json:"notes" gorm:"type:text"`

	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index:idx_appointments_deleted_at"`
}
