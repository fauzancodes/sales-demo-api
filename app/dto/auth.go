package dto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (request LoginRequest) Validate() error {
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
	Email                  string `json:"email"`
	SuccessVerificationUrl string `json:"success_verification_url"`
	FailedVerificationUrl  string `json:"failed_verification_url"`
}

func (request ResendEmailVerification) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Email, validation.Required, is.Email),
		validation.Field(&request.SuccessVerificationUrl, is.URL),
		validation.Field(&request.FailedVerificationUrl, is.URL),
	)
}

type RegisterRequest struct {
	Email                  string `json:"email"`
	Password               string `json:"password"`
	SuccessVerificationUrl string `json:"success_verification_url"`
	FailedVerificationUrl  string `json:"failed_verification_url"`
}

func (request RegisterRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Email, validation.Required, is.Email),
		validation.Field(&request.Password, validation.Required),
		validation.Field(&request.SuccessVerificationUrl, is.URL),
		validation.Field(&request.FailedVerificationUrl, is.URL),
	)
}
