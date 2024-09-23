package models

type SDAUser struct {
	CustomGormModel
	FirstName string `json:"first_name" gorm:"type:varchar(255);column:first_name"`
	LastName  string `json:"last_name" gorm:"type:varchar(255);column:last_name"`
	Email     string `json:"email" gorm:"type:varchar(255);column:email"`
	Password  string `json:"password" gorm:"type:varchar(255);column:password"`
}

func (SDAUser) TableName() string {
	return "sda_users"
}
