package models

import "github.com/google/uuid"

type SDAProductCategory struct {
	CustomGormModel
	Name        string    `json:"name" gorm:"type:varchar(255);column:name"`
	Description string    `json:"description" gorm:"type:text;column:description"`
	Status      bool      `json:"status" gorm:"type:bool;column:status"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;column:user_id"`
}

func (SDAProductCategory) TableName() string {
	return "sda_product_categories"
}
