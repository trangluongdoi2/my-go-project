package serviceoffering

import (
	"time"

	"github.com/google/uuid"
)

type Staff struct {
	ID        uuid.UUID `gorm:"type:uuid"`
	Code      string
	FirstName string
	LastName  string
	Email     string
	Phone     string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (Staff) TableName() string {
	return "staffs"
}

type ServiceOffering struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Code        string    `json:"code" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Duration    int       `json:"duration" gorm:"comment:Duration in minutes"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Status      int16     `json:"status" gorm:"type:smallint;default:1"`

	Staffs *[]Staff `json:"staffs,omitempty" gorm:"many2many:staff_services;foreignKey:ID;joinForeignKey:ServiceOfferingID;References:ID;joinReferences:StaffID"`

	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index:idx_service_offerings_deleted_at"`
}

type StaffService struct {
	StaffID           uuid.UUID `json:"staff_id" gorm:"type:uuid;primaryKey"`
	ServiceOfferingID uuid.UUID `json:"service_offering_id" gorm:"type:uuid;primaryKey"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
}
