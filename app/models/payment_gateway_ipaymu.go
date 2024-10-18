package models

type SDAIPaymuPaymentMethod struct {
	CustomGormModel
	Code        string `json:"code" gorm:"type:varchar(50);column:code"`
	Name        string `json:"name" gorm:"type:varchar(50);column:name"`
	Description string `json:"decription" gorm:"type:text;column:decription"`
}

func (SDAIPaymuPaymentMethod) TableName() string {
	return "sda_ipaymu_payment_methods"
}
