package models

import "github.com/google/uuid"

type SDAProduct struct {
	CustomGormModel
	Code        string                  `json:"code" gorm:"type:varchar(50);column:code"`
	Name        string                  `json:"name" gorm:"type:varchar(255);column:name"`
	Description string                  `json:"description" gorm:"type:text;column:description"`
	Image       string                  `json:"image" gorm:"type:text;column:image"`
	Status      bool                    `json:"status" gorm:"type:bool;column:status"`
	Price       float64                 `json:"price" gorm:"type:float8;column:price"`
	CategoryID  uuid.UUID               `json:"category_id" gorm:"type:uuid;column:category_id"`
	UserID      uuid.UUID               `json:"user_id" gorm:"type:uuid;column:user_id"`
	Category    ProductCategoryRelation `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	User        UserRelation            `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Stocks      []ProductStockRelation  `json:"stocks,omitempty" gorm:"foreignKey:ProductID"`
}

func (SDAProduct) TableName() string {
	return "sda_products"
}

type ProductRelation struct {
	CustomGormModel
	Code       string    `json:"code" gorm:"column:code"`
	Name       string    `json:"name" gorm:"column:name"`
	Status     bool      `json:"status" gorm:"column:status"`
	Price      float64   `json:"price" gorm:"column:price"`
	UserID     uuid.UUID `json:"-" gorm:"column:user_id"`
	CategoryID uuid.UUID `json:"-" gorm:"column:category_id"`
}

func (ProductRelation) TableName() string {
	return "sda_products"
}
