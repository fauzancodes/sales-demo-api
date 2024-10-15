package dto

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type MidtransRequest struct {
	InvoiceID         string              `json:"invoice_id"`
	PaymentMethodCode string              `json:"payment_method_code"`
	Card              MidtransCardRequest `json:"card"`
}

func (request MidtransRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.InvoiceID, validation.Required),
		validation.Field(&request.PaymentMethodCode, validation.Required),
	)
}

type MidtransCardRequest struct {
	CardNumber string `json:"card_number"`
	ExpMonth   int    `json:"exp_month"`
	ExpYear    int    `json:"exp_year"`
	CVV        string `json:"cvv"`
}

func (request MidtransCardRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.CardNumber, validation.Required),
		validation.Field(&request.ExpMonth, validation.Required, validation.Min(1), validation.Max(12)),
		validation.Field(&request.ExpYear, validation.Required, validation.Min(2000)),
		validation.Field(&request.CVV, validation.Required),
	)
}
