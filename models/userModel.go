package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	FirstName    string     `gorm:"not null;type:varchar(100)" json:"first_name"`
	LastName     string     `gorm:"not null;type:varchar(100)" json:"last_name"`
	Password     string     `gorm:"not null;type:varchar(100)" json:"-"`
	Email        string     `gorm:"unique;not null;type:varchar(100)" json:"email"`
	Phone        *string    `gorm:"unique;type:varchar(16);default:null" json:"phone"`
	BirthDay     *time.Time `json:"birth_day"`
	AvatarUrl    *string    `json:"avatar_url"`
	Website      *string    `json:"website"`
	Token        *string    `json:"token"`
	Role         string     `gorm:"not null;type:varchar(16);default:USER" json:"role"`
	RefreshToken *string    `gorm:"type:text" json:"refresh_token"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoCreateTime" json:"updated_at"`
	UserId       string     `gorm:"autoUpdateTime" json:"user_id"`
}
