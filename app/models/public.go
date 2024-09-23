package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomGormModel struct {
	ID        uuid.UUID       `gorm:"type:uuid;column:id;default:uuid_generate_v4();primaryKey" json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"-"`
}
