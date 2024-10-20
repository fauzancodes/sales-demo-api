package models

type SDAXenditPaymentMethod struct {
	CustomGormModel
	Code        string `json:"code" gorm:"type:varchar(50);column:code"`
	Name        string `json:"name" gorm:"type:varchar(50);column:name"`
	Description string `json:"decription" gorm:"type:text;column:decription"`
}

func (SDAXenditPaymentMethod) TableName() string {
	return "sda_xendit_payment_methods"
}
