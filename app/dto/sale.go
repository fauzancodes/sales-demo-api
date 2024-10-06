package dto

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type SaleRequest struct {
	InvoiceID       string              `json:"invoice_id"`
	Discount        float64             `json:"discount"`
	Tax             float64             `json:"tax"`
	MiscPrice       float64             `json:"misc_price"`
	Subtotal        float64             `json:"subtotal"`
	TotalPaid       float64             `json:"total_paid"`
	CustomerID      string              `json:"customer_id"`
	TransactionDate string              `json:"transaction_date"`
	Details         []SaleDetailRequest `json:"details"`
}

func (request SaleRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.TotalPaid, validation.Min(0.0)),
		validation.Field(&request.Subtotal, validation.Min(0.0)),
		validation.Field(&request.MiscPrice, validation.Min(0.0)),
		validation.Field(&request.Discount, validation.Min(0.0), validation.Max(100.0)),
		validation.Field(&request.Tax, validation.Min(0.0), validation.Max(100.0)),
	)
}

type SaleDetailRequest struct {
	ProductID string  `json:"product_id"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}

func (request SaleDetailRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.ProductID, validation.Required),
		validation.Field(&request.Quantity, validation.Required),
		validation.Field(&request.Price, validation.Min(0.0)),
	)
}
