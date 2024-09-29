package models

import "github.com/google/uuid"

type SDACustomer struct {
	CustomGormModel
	FirstName string    `json:"first_name" gorm:"type:varchar(255);column:first_name"`
	LastName  string    `json:"last_name" gorm:"type:varchar(255);column:last_name"`
	Email     string    `json:"email" gorm:"type:varchar(255);column:email"`
	Phone     string    `json:"phone" gorm:"type:varchar(255);column:phone"`
	Status    bool      `json:"status" gorm:"type:bool;column:status"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;column:user_id"`
}

func (SDACustomer) TableName() string {
	return "sda_customers"
}
