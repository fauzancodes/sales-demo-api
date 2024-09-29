package dto

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type ProductStockRequest struct {
	Amount      int    `json:"amount"`
	Description string `json:"description"`
	ProductID   string `json:"product_id"`
	Action      string `json:"action"`
}

func (request ProductStockRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Amount, validation.Required),
		validation.Field(&request.ProductID, validation.Required),
		validation.Field(&request.Action, validation.Required, validation.In("add", "reduce")),
	)
}
