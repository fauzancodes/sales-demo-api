package models

import "github.com/google/uuid"

type SDASaleDetail struct {
	CustomGormModel
	ProductID uuid.UUID        `json:"product_id" gorm:"type:uuid;column:product_id"`
	Price     float64          `json:"price" gorm:"type:float8;column:price"`
	Quantity  int              `json:"quantity" gorm:"type:int8;column:quantity"`
	SaleID    uuid.UUID        `json:"sale_id" gorm:"type:uuid;column:sale_id"`
	UserID    uuid.UUID        `json:"user_id" gorm:"type:uuid;column:user_id"`
	Sale      *SaleRelation    `json:"sale,omitempty" gorm:"foreignKey:SaleID"`
	User      *UserRelation    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Product   *ProductRelation `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

func (SDASaleDetail) TableName() string {
	return "sda_sale_details"
}

type SaleDetailRelation struct {
	CustomGormModel
	ProductID uuid.UUID        `json:"product_id" gorm:"type:uuid;column:product_id"`
	Price     float64          `json:"price" gorm:"type:float8;column:price"`
	Quantity  int              `json:"quantity" gorm:"type:int8;column:quantity"`
	UserID    uuid.UUID        `json:"-" gorm:"type:uuid;column:user_id"`
	SaleID    uuid.UUID        `json:"-" gorm:"type:uuid;column:sale_id"`
	Product   *ProductRelation `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

func (SaleDetailRelation) TableName() string {
	return "sda_sale_details"
}
