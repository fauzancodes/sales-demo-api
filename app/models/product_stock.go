package models

import "github.com/google/uuid"

type SDAProductStock struct {
	CustomGormModel
	ProductID   uuid.UUID `json:"product_id" gorm:"type:uuid;column:product_id"`
	Addition    int       `json:"addition" gorm:"type:int8;column:addition"`
	Reduction   int       `json:"reduction" gorm:"type:int8;column:reduction"`
	Current     int       `json:"current" gorm:"type:int8;column:current"`
	Description string    `json:"description" gorm:"type:text;column:description"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;column:user_id"`
}

func (SDAProductStock) TableName() string {
	return "sda_product_stocks"
}
