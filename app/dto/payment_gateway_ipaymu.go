package dto

import validation "github.com/go-ozzo/ozzo-validation"

type IPaymuSaleRequest struct {
	InvoiceID         string `json:"invoice_id"`
	PaymentMethodCode string `json:"payment_method_code"`
}

func (request IPaymuSaleRequest) Validate() error {
	return validation.ValidateStruct(
		&request,
		validation.Field(&request.InvoiceID, validation.Required),
		validation.Field(&request.PaymentMethodCode, validation.Required),
	)
}

type IPaymuRequest struct {
	Name            string   `json:"name"`
	Phone           string   `json:"phone"`
	Email           string   `json:"email"`
	Amount          int64    `json:"amount"`
	NotifyURL       string   `json:"notifyUrl"`
	Expired         int      `json:"expired"`
	ReferenceID     string   `json:"referenceId"`
	PaymentMethod   string   `json:"paymentMethod"`
	PaymentChannel  string   `json:"paymentChannel"`
	ProductName     []string `json:"product"`
	ProductQuantity []int    `json:"qty"`
	ProductPrice    []int64  `json:"price"`
}

type IPaymuResponseData struct {
	SessionID      string  `json:"SessionId"`
	TransactionID  int     `json:"TransactionId"`
	ReferenceID    string  `json:"ReferenceId"`
	PaymentMethod  string  `json:"Via"`
	PaymentChannel string  `json:"Channel"`
	PaymentCode    string  `json:"PaymentNo"`
	TargetName     string  `json:"PaymentName"`
	Subtotal       float64 `json:"Total"`
	Fee            float64 `json:"Fee"`
	Expired        string  `json:"Expired"`
}

type IPaymuResponse struct {
	Status  int                `json:"Status"`
	Message string             `json:"Message"`
	Data    IPaymuResponseData `json:"Data"`
}

type IPaymuNotificationRequest struct {
	TransactionID int64  `json:"trx_id"`
	Status        string `json:"status"`
	StatusCode    string `json:"status_code"`
	SID           string `json:"sid"`
	ReferenceID   string `json:"reference_id"`
}
