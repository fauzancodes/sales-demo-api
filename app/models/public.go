package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomGormModel struct {
	ID        uuid.UUID       `gorm:"type:uuid;column:id;default:uuid_generate_v4();primaryKey" json:"id"`
	CreatedAt time.Time       `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"-"`
}

type SDAUsedApiKey struct {
	CustomGormModel
	SecretKey    string `json:"-" gorm:"type: text"`
	Base64Key    string `json:"-" gorm:"type: text"`
	ReceivedHMAC string `json:"-" gorm:"type: text"`
	ExpectedHMAC string `json:"-" gorm:"type: text"`
}

func (SDAUsedApiKey) TableName() string {
	return "sda_used_api_keys"
}
