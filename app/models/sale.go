package models

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type SDASale struct {
	CustomGormModel
	InvoiceID       string               `json:"invoice_id" gorm:"type:varchar(255);column:invoice_id"`
	CustomerID      uuid.UUID            `json:"customer_id" gorm:"type:uuid;column:customer_id"`
	Discount        float64              `json:"discount" gorm:"type:float8;column:discount"`
	Tax             float64              `json:"tax" gorm:"type:float8;column:tax"`
	MiscPrice       float64              `json:"misc_price" gorm:"type:float8;column:misc_price"`
	Subtotal        float64              `json:"subtotal" gorm:"type:float8;column:subtotal"`
	TotalPaid       float64              `json:"total_paid" gorm:"type:float8;column:total_paid"`
	Status          bool                 `json:"status" gorm:"type:bool;column:status"`
	TransactionDate null.Time            `json:"transaction_date" gorm:"type:timestamptz;column:transaction_date"`
	UserID          uuid.UUID            `json:"user_id" gorm:"type:uuid;column:user_id"`
	User            UserRelation         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Customer        CustomerRelation     `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Details         []SaleDetailRelation `json:"details,omitempty" gorm:"foreignKey:SaleID"`
}

func (SDASale) TableName() string {
	return "sda_sales"
}

type SaleRelation struct {
	CustomGormModel
	InvoiceID       string    `json:"invoice_id" gorm:"column:invoice_id"`
	Discount        float64   `json:"discount" gorm:"column:discount"`
	Tax             float64   `json:"tax" gorm:"column:tax"`
	MiscPrice       float64   `json:"misc_price" gorm:"column:misc_price"`
	Subtotal        float64   `json:"subtotal" gorm:"column:subtotal"`
	TotalPaid       float64   `json:"total_paid" gorm:"column:total_paid"`
	TransactionDate null.Time `json:"transaction_date" gorm:"column:transaction_date"`
	UserID          uuid.UUID `json:"-" gorm:"type:uuid;column:user_id"`
	CustomerID      uuid.UUID `json:"-" gorm:"type:uuid;column:customer_id"`
}

func (SaleRelation) TableName() string {
	return "sda_sales"
}
