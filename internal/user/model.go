package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Code      string    `json:"code"`
	FirstName string    `json:"first_name" gorm:"not null"`
	LastName  string    `json:"last_name" gorm:"not null"`
	Email     string    `json:"email" gorm:"uniqueIndex:idx_users_email"`
	Phone     string    `json:"phone"`
	Password  string    `json:"password" gorm:"not null"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	Role      string    `json:"role" gorm:"default:'customer'"`

	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index:idx_users_deleted_at"`
}
