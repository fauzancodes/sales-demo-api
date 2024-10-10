package dto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type CustomerRequest struct {
	Code      string `json:"code"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Status    bool   `json:"status"`
}

func (request CustomerRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.FirstName, validation.Required),
		validation.Field(&request.Email, validation.Required, is.Email),
	)
}
