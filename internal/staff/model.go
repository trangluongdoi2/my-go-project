package staff

import (
	baserepo "go-backend-project/internal/repository"
	"time"

	"github.com/google/uuid"
)

type Staff struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Code      string    `json:"code"`
	FirstName string    `json:"first_name" gorm:"not null"`
	LastName  string    `json:"last_name" gorm:"not null"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone" gorm:"index:idx_staffs_phone,unique"`
	IsActive  bool      `json:"is_active"`

	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type ListStaffQuery struct {
	baserepo.BaseQuery

	Code      string `form:"code"`
	FirstName string `form:"first_name"`
	LastName  string `form:"last_name"`
	Email     string `form:"email"`
	Phone     string `form:"phone"`
	IsActive  *bool  `form:"is_active"`
}
