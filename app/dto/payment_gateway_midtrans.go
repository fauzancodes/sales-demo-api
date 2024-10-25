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
	if request.PaymentMethodCode == "credit_card" {
		err := request.Card.Validate()
		if err != nil {
			return err
		}
	}

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

type MidtransNotificationRequest struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	OrderID           string `json:"order_id"`
	SignatureKey      string `json:"signature_key"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
}
