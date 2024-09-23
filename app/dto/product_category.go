package dto

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type ProductCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
}

func (request ProductCategoryRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Name, validation.Required),
	)
}
