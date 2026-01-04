package staff

import (
	"time"

	"github.com/google/uuid"
)

type Staff struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Code      string    `json:"code"`
	FirstName string    `json:"first_name" gorm:"not null"`
	LastName  string    `json:"last_name" gorm:"not null"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	IsActive  bool      `json:"is_active"`

	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at"`
}
