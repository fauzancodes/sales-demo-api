package models

type SDAUser struct {
	CustomGormModel
	FirstName         string                    `json:"first_name" gorm:"type:varchar(255);column:first_name"`
	LastName          string                    `json:"last_name" gorm:"type:varchar(255);column:last_name"`
	Email             string                    `json:"email" gorm:"type:varchar(255);column:email"`
	Password          string                    `json:"-" gorm:"type:varchar(255);column:password"`
	Products          []ProductRelation         `json:"products,omitempty" gorm:"foreignKey:UserID"`
	ProductCategories []ProductCategoryRelation `json:"product_categories,omitempty" gorm:"foreignKey:UserID"`
	ProductStocks     []ProductStockRelation    `json:"product_stocks,omitempty" gorm:"foreignKey:UserID"`
	Customers         []CustomerRelation        `json:"customers,omitempty" gorm:"foreignKey:UserID"`
	Sales             []SaleRelation            `json:"sales,omitempty" gorm:"foreignKey:UserID"`
	SaleDetails       []SaleDetailRelation      `json:"sale_details,omitempty" gorm:"foreignKey:UserID"`
}

func (SDAUser) TableName() string {
	return "sda_users"
}

type UserRelation struct {
	CustomGormModel
	FirstName string `json:"first_name" gorm:"column:first_name"`
	LastName  string `json:"last_name" gorm:"column:last_name"`
	Email     string `json:"email" gorm:"type:column:email"`
}

func (UserRelation) TableName() string {
	return "sda_users"
}
