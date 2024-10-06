package models

import "github.com/google/uuid"

type SDACustomer struct {
	CustomGormModel
	FirstName string         `json:"first_name" gorm:"type:varchar(255);column:first_name"`
	LastName  string         `json:"last_name" gorm:"type:varchar(255);column:last_name"`
	Email     string         `json:"email" gorm:"type:varchar(255);column:email"`
	Phone     string         `json:"phone" gorm:"type:varchar(255);column:phone"`
	Status    bool           `json:"status" gorm:"type:bool;column:status"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;column:user_id"`
	User      UserRelation   `gorm:"foreignKey:UserID"`
	Sales     []SaleRelation `gorm:"foreignKey:CustomerID"`
}

func (SDACustomer) TableName() string {
	return "sda_customers"
}

type CustomerRelation struct {
	CustomGormModel
	FirstName string    `json:"first_name" gorm:"column:first_name"`
	LastName  string    `json:"last_name" gorm:"column:last_name"`
	Email     string    `json:"email" gorm:"column:email"`
	Phone     string    `json:"phone" gorm:"column:phone"`
	Status    bool      `json:"status" gorm:"column:status"`
	UserID    uuid.UUID `json:"-" gorm:"column:user_id"`
}

func (CustomerRelation) TableName() string {
	return "sda_customers"
}
