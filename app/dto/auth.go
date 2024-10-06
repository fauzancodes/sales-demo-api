package dto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (request AuthRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Email, validation.Required, is.Email),
		validation.Field(&request.Password, validation.Required),
	)
}

type EmailVerfication struct {
	Name            string
	VerificationUrl string
	AppUrl          string
}

type ResendEmailVerification struct {
	Email string `json:"email"`
}
