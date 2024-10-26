package models

import (
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type SDAMidtransPaymentMethod struct {
	CustomGormModel
	Code        string `json:"code" gorm:"type:varchar(50);column:code"`
	Name        string `json:"name" gorm:"type:varchar(50);column:name"`
	Description string `json:"decription" gorm:"type:text;column:decription"`
}

func (SDAMidtransPaymentMethod) TableName() string {
	return "sda_midtrans_payment_methods"
}

type SDAMidtransSalePayment struct {
	CustomGormModel
	SaleID            uuid.UUID                 `json:"sale_id" gorm:"type:uuid;column:sale_id"`
	ReferenceCode     string                    `json:"reference_code" gorm:"type:varchar(255);column:reference_code"`
	ExpiryDate        null.Time                 `json:"expiry_date" gorm:"type:timestamptz;column:expiry_date"`
	RawResponse       string                    `json:"-" gorm:"type:text;column:raw_response"`
	PaymentCode       string                    `json:"payment_code" gorm:"type:varchar(50);column:payment_code"`
	MandiriBillKey    string                    `json:"mandiri_bill_key" gorm:"type:varchar(50);column:mandiri_bill_key"`
	MandiriBillerCode string                    `json:"mandiri_biller_code" gorm:"type:varchar(50);column:mandiri_biller_code"`
	QRCodeUrl         string                    `json:"qr_code_url" gorm:"type:varchar(255);column:qr_code_url"`
	RedirectUrl       string                    `json:"redirect_url" gorm:"type:varchar(255);column:redirect_url"`
	UserID            uuid.UUID                 `json:"user_id" gorm:"type:uuid;column:user_id"`
	MerchantID        string                    `json:"merchant_id" gorm:"type:varchar(50);column:merchant_id"`
	PaymentMethodID   *uuid.UUID                `json:"payment_method_id" gorm:"type:uuid;column:payment_method_id"`
	PaymentMethod     *SDAMidtransPaymentMethod `json:"payment_method,omitempty" gorm:"foreignKey:PaymentMethodID"`
}

func (SDAMidtransSalePayment) TableName() string {
	return "sda_midtrans_sale_payments"
}
