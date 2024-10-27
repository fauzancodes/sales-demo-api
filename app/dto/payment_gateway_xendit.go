package dto

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type XenditRequestPayment struct {
	InvoiceID         string                   `json:"invoice_id"`
	PaymentMethodCode string                   `json:"payment_method_code"`
	SuccessReturnUrl  string                   `json:"success_return_url"`
	FailedReturnUrl   string                   `json:"failed_return_url"`
	Card              XenditCardRequest        `json:"card"`
	EWallet           XenditEWalletRequest     `json:"ewallet"`
	DirectDebit       XenditDirectDebitRequest `json:"direct_debit"`
}

func (request XenditRequestPayment) Validate() error {
	var err error
	if request.PaymentMethodCode == "DIRECT_DEBIT_BRI" || request.PaymentMethodCode == "DIRECT_DEBIT_MANDIRI" {
		if request.PaymentMethodCode == "DIRECT_DEBIT_BRI" {
			if request.DirectDebit.MobileNumber == "" {
				err = errors.New("direct_debit.mobile_number is required for BRI")
				return err
			}
			if request.DirectDebit.LastFourDigit == "" {
				err = errors.New("direct_debit.last_four_digit is required for BRI")
				return err
			}
		}

		if request.DirectDebit.Type == "" {
			err = errors.New("direct_debit.type is required and only accepts BANK_ACCOUNT or DEBIT_CARD")
			return err
		}

		if request.DirectDebit.Type != "BANK_ACCOUNT" && request.DirectDebit.Type != "DEBIT_CARD" {
			err = errors.New("direct_debit.type only accepts BANK_ACCOUNT or DEBIT_CARD")
			return err
		}

		err = request.DirectDebit.Validate()
		if err != nil {
			return err
		}
	}

	if request.PaymentMethodCode == "OVO" {
		if request.EWallet.MobileNumber == "" {
			err = errors.New("ewallet.mobile_number is required for OVO")
			return err
		}

		var customerPhone string
		if request.EWallet.MobileNumber[:2] == "08" {
			customerPhone = "+62" + request.EWallet.MobileNumber[1:]
		} else if request.EWallet.MobileNumber[:3] == "+62" {
			customerPhone = request.EWallet.MobileNumber
		} else {
			customerPhone = "+62" + request.EWallet.MobileNumber
		}

		request.EWallet.MobileNumber = customerPhone
	}

	if request.PaymentMethodCode == "JENIUSPAY" {
		if request.EWallet.CashTag == "" {
			err = errors.New("ewallet.cashtag is required for JENIUSPAY")
			return err
		}
	}

	return validation.ValidateStruct(
		&request,
		validation.Field(&request.InvoiceID, validation.Required),
		validation.Field(&request.PaymentMethodCode, validation.Required),
		validation.Field(&request.SuccessReturnUrl, validation.Required),
		validation.Field(&request.FailedReturnUrl, validation.Required),
	)
}

type XenditRequestInvoice struct {
	InvoiceID string `json:"invoice_id"`
}

func (request XenditRequestInvoice) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.InvoiceID, validation.Required),
	)
}

type XenditEWalletRequest struct {
	MobileNumber string `json:"mobile_number"`
	CashTag      string `json:"cashtag"`
}

type XenditDirectDebitRequest struct {
	Type          string `json:"type"`
	MobileNumber  string `json:"mobile_number"`
	LastFourDigit string `json:"last_four_digit"`
	CardExpiry    string `json:"card_expiry"`
	Email         string `json:"email"`
}

func (request XenditDirectDebitRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.Email, is.Email),
		validation.Field(&request.Type, validation.In("BANK_ACCOUNT", "DEBIT_CARD")),
	)
}

type XenditCardRequest struct {
	CardNumber            string `json:"card_number"`
	ExpiryMonth           string `json:"expiry_month"`
	ExpiryYear            string `json:"expiry_year"`
	CardHolderName        string `json:"cardholder_name"`
	CardHolderEmail       string `json:"cardholder_email"`
	CardHolderPhoneNumber string `json:"cardholder_phone_number"`
}

func (request XenditCardRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.CardNumber, validation.Required),
		validation.Field(&request.ExpiryMonth, validation.Required),
		validation.Field(&request.ExpiryYear, validation.Required),
		validation.Field(&request.CardHolderName, validation.Required),
		validation.Field(&request.CardHolderEmail, validation.Required, is.Email),
		validation.Field(&request.CardHolderPhoneNumber, validation.Required),
	)
}

type XenditNotificationRequest struct {
	Event      string                        `json:"event"`
	BusinessID string                        `json:"business_id"`
	Created    string                        `json:"created"`
	Data       XenditNotificationDataRequest `json:"data"`
}

type XenditNotificationDataRequest struct {
	ReferenceID string `json:"reference_id"`
}
