package service

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/repository"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

func GetMidtransPaymentMethods(code string, param utils.PagingRequest) (response utils.PagingResponse, data []models.SDAMidtransPaymentMethod, statusCode int, err error) {
	baseFilter := "deleted_at IS NULL"
	filter := baseFilter
	var filterValues []any

	if code != "" {
		filter += " AND code = ?"
		filterValues = append(filterValues, code)
	}
	if param.Search != "" {
		filter += " AND (name ILIKE ? OR description ILIKE ?)"
		filterValues = append(filterValues, fmt.Sprintf("%%%s%%", param.Search))
		filterValues = append(filterValues, fmt.Sprintf("%%%s%%", param.Search))
	}

	data, total, totalFiltered, err := repository.GetMidtransPaymentMethods(dto.FindParameter{
		BaseFilter:   baseFilter,
		Filter:       filter,
		FilterValues: filterValues,
		Limit:        param.Limit,
		Order:        param.Order,
		Offset:       param.Offset,
	})
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	response = utils.PopulateResPaging(&param, data, total, totalFiltered)

	statusCode = http.StatusOK
	return
}

func MidtransChargeCore(userID, baseUrl string, request dto.MidtransRequestCore) (response models.SDAMidtransSalePayment, statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = errors.New("failed to parse user UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	paymentMethodData, _, _, err := repository.GetMidtransPaymentMethods(dto.FindParameter{
		Filter:       "deleted_at IS NULL AND code = ?",
		FilterValues: []any{strings.ToLower(request.PaymentMethodCode)},
	})
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}
	if len(paymentMethodData) == 0 {
		err = errors.New("payment method not found")
		statusCode = http.StatusNotFound
		return
	}
	paymentMethod := paymentMethodData[0]

	saleData, _, _, err := repository.GetSales(dto.FindParameter{
		Filter:       "deleted_at IS NULL AND invoice_id = ?",
		FilterValues: []any{request.InvoiceID},
	}, []string{"Details", "Details.Product", "Customer"})
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}
	if len(saleData) == 0 {
		err = errors.New("sale data not found")
		statusCode = http.StatusNotFound
		return
	}
	if len(saleData[0].Details) == 0 {
		err = errors.New("sale details data not found")
		statusCode = http.StatusNotFound
		return
	}

	sale := saleData[0]

	serverKey := config.LoadConfig().MidtransServerKey
	env := midtrans.Sandbox
	if strings.ToLower(config.LoadConfig().MidtransEnv) == "production" {
		env = midtrans.Production
	}

	c := coreapi.Client{}
	c.New(serverKey, env)
	c.Options.SetPaymentOverrideNotification(baseUrl + "/payment-gateway/midtrans/notification")

	chargeReq := &coreapi.ChargeReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  sale.InvoiceID,
			GrossAmt: int64(utils.RoundFloat(sale.TotalPaid)),
		},
		CustomerDetails: &midtrans.CustomerDetails{
			FName: sale.Customer.FirstName,
			LName: sale.Customer.LastName,
			Email: sale.Customer.Email,
			Phone: sale.Customer.Phone,
		},
		CustomExpiry: &coreapi.CustomExpiry{
			OrderTime:      time.Now().Format("2006-01-02 15:04:05 +0700"),
			ExpiryDuration: 1,
			Unit:           "day",
		},
	}

	switch paymentMethod.Code {
	case "credit_card":
		cardResponse, midtransError := c.CardToken(strings.ReplaceAll(request.Card.CardNumber, " ", ""), request.Card.ExpMonth, request.Card.ExpYear, request.Card.CVV, config.LoadConfig().MidtransClientKey)
		if midtransError != nil {
			err = errors.New("failed to get card token: " + midtransError.RawError.Error())
			statusCode = http.StatusInternalServerError
			return
		}

		chargeReq.PaymentType = coreapi.PaymentTypeCreditCard
		chargeReq.CreditCard = &coreapi.CreditCardDetails{
			TokenID:        cardResponse.TokenID,
			Authentication: true,
		}

		var bank string
		switch cardResponse.Bank {
		case "bca":
			bank = string(midtrans.BankBca)
		case "mandiri":
			bank = string(midtrans.BankMandiri)
		case "bni":
			bank = string(midtrans.BankBni)
		case "cimb":
			bank = string(midtrans.BankCimb)
		case "maybank":
			bank = string(midtrans.BankMaybank)
		case "bri":
			bank = string(midtrans.BankBri)
		}
		chargeReq.CreditCard = &coreapi.CreditCardDetails{Bank: bank}
	case "akulaku":
		chargeReq.PaymentType = coreapi.PaymentTypeAkulaku
	case "kredivo":
		chargeReq.PaymentType = coreapi.CoreapiPaymentType(paymentMethod.Code)
	case "qris_gopay":
		chargeReq.PaymentType = coreapi.PaymentTypeQris
		chargeReq.Qris = &coreapi.QrisDetails{
			Acquirer: "gopay",
		}
	case "qris_shopeepay":
		chargeReq.PaymentType = coreapi.PaymentTypeQris
		chargeReq.Qris = &coreapi.QrisDetails{
			Acquirer: "airpay shopee",
		}
	case "gopay":
		chargeReq.PaymentType = coreapi.PaymentTypeGopay
		chargeReq.Gopay = &coreapi.GopayDetails{
			EnableCallback: true,
			CallbackUrl:    baseUrl,
		}
	case "shopeepay":
		chargeReq.PaymentType = coreapi.PaymentTypeShopeepay
		chargeReq.ShopeePay = &coreapi.ShopeePayDetails{
			CallbackUrl: baseUrl,
		}
		chargeReq.CustomExpiry = &coreapi.CustomExpiry{
			OrderTime:      time.Now().Format("2006-01-02 15:04:05 +0700"),
			ExpiryDuration: 1,
			Unit:           "day",
		}
	case "alfamart":
		chargeReq.PaymentType = coreapi.PaymentTypeConvenienceStore
		chargeReq.ConvStore = &coreapi.ConvStoreDetails{
			Store:             "alfamart",
			Message:           "Sale Payment",
			AlfamartFreeText1: "Invoice ID: " + sale.InvoiceID,
			AlfamartFreeText2: "Generate By: Sales Demo API",
		}
	case "indomaret":
		chargeReq.PaymentType = coreapi.PaymentTypeConvenienceStore
		chargeReq.ConvStore = &coreapi.ConvStoreDetails{
			Store:   "indomaret",
			Message: "Sale Payment. Invoice ID: " + sale.InvoiceID + " Generate By: Sales Demo API",
		}
	case "mandiri":
		chargeReq.PaymentType = coreapi.PaymentTypeEChannel
		chargeReq.EChannel = &coreapi.EChannelDetail{
			BillInfo1: "Sale Payment",
			BillInfo2: "Invoice ID: " + sale.InvoiceID,
			BillInfo3: "Generate By: Sales Demo API",
		}
	case "permata":
		chargeReq.PaymentType = coreapi.CoreapiPaymentType("permata")
	default:
		chargeReq.PaymentType = coreapi.PaymentTypeBankTransfer
		var bank midtrans.Bank
		switch paymentMethod.Code {
		case "bca":
			bank = midtrans.BankBca
		case "bni":
			bank = midtrans.BankBni
		case "bri":
			bank = midtrans.BankBri
		case "cimb":
			bank = midtrans.BankCimb
		}
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{Bank: bank}
	}

	var totalQuantity int
	for _, item := range sale.Details {
		totalQuantity += item.Quantity
	}

	var items []midtrans.ItemDetails
	for _, data := range sale.Details {
		if sale.Discount > 0 {
			discountAmount := (sale.Discount * sale.Subtotal) / 100
			discountPerItem := discountAmount / float64(totalQuantity)
			if discountPerItem > 0 {
				data.Price -= discountPerItem
			}
		}
		if sale.Tax > 0 {
			taxAmount := (sale.Tax * sale.Subtotal) / 100
			taxPerItem := taxAmount / float64(totalQuantity)
			if taxPerItem > 0 {
				data.Price += taxPerItem
			}
		}
		if sale.MiscPrice > 0 {
			miscPricePerItem := sale.MiscPrice / float64(totalQuantity)
			if miscPricePerItem > 0 {
				data.Price += miscPricePerItem
			}
		}

		item := midtrans.ItemDetails{
			ID:    data.ProductID.String(),
			Name:  data.Product.Name,
			Qty:   int32(data.Quantity),
			Price: int64(utils.RoundFloat(data.Price)),
		}

		items = append(items, item)
	}
	if len(items) > 0 {
		chargeReq.Items = &items
	}

	midtransResponse, midtransError := c.ChargeTransaction(chargeReq)
	if midtransError != nil {
		err = errors.New("failed to charge to midtrans: " + midtransError.RawError.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		err = errors.New("failed to get location time: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	rawResponse, err := json.Marshal(midtransResponse)
	if err != nil {
		err = errors.New("failed to marshal response: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data := models.SDAMidtransSalePayment{
		SaleID:          sale.ID,
		PaymentMethodID: &paymentMethod.ID,
		ReferenceCode:   midtransResponse.TransactionID,
		ExpiryDate:      null.TimeFrom(time.Now().In(location).Add(24 * time.Hour)),
		RawResponse:     string(rawResponse),
		UserID:          parsedUserUUID,
		MerchantID:      config.LoadConfig().MidtransMerchantID,
	}

	switch paymentMethod.Code {
	case "credit_card", "akulaku", "kredivo":
		data.RedirectUrl = midtransResponse.RedirectURL
	case "qris_gopay":
		if len(midtransResponse.Actions) > 0 {
			data.QRCodeUrl = midtransResponse.Actions[0].URL
		}
	case "qris_shopeepay":
		if len(midtransResponse.Actions) > 0 {
			data.QRCodeUrl = midtransResponse.Actions[0].URL
		}
		data.ExpiryDate = null.TimeFrom(time.Now().In(location).Add(time.Hour))
	case "gopay":
		if len(midtransResponse.Actions) > 0 {
			data.QRCodeUrl = midtransResponse.Actions[0].URL
			data.RedirectUrl = midtransResponse.Actions[1].URL
		}
	case "shopeepay":
		if len(midtransResponse.Actions) > 0 {
			data.RedirectUrl = midtransResponse.Actions[0].URL
		}
		data.ExpiryDate = null.TimeFrom(time.Now().In(location).Add(time.Hour))
	case "alfamart", "indomaret":
		data.PaymentCode = midtransResponse.PaymentCode
	case "mandiri":
		data.MandiriBillKey = midtransResponse.BillKey
		data.MandiriBillerCode = midtransResponse.BillerCode
	case "permata":
		data.PaymentCode = midtransResponse.PermataVaNumber
	default:
		if len(midtransResponse.VaNumbers) > 0 {
			data.PaymentCode = midtransResponse.VaNumbers[0].VANumber
		}
	}

	response, err = repository.CreateMidtransSalePayment(data)
	if err != nil {
		err = errors.New("failed to create data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusCreated
	return
}

func MidtransChargeSnap(userID, baseUrl string, request dto.MidtransRequestSnap) (response models.SDAMidtransSalePayment, statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = errors.New("failed to parse user UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	saleData, _, _, err := repository.GetSales(dto.FindParameter{
		Filter:       "deleted_at IS NULL AND invoice_id = ?",
		FilterValues: []any{request.InvoiceID},
	}, []string{"Details", "Details.Product", "Customer"})
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}
	if len(saleData) == 0 {
		err = errors.New("sale data not found")
		statusCode = http.StatusNotFound
		return
	}
	if len(saleData[0].Details) == 0 {
		err = errors.New("sale details data not found")
		statusCode = http.StatusNotFound
		return
	}

	sale := saleData[0]

	serverKey := config.LoadConfig().MidtransServerKey
	env := midtrans.Sandbox
	if strings.ToLower(config.LoadConfig().MidtransEnv) == "production" {
		env = midtrans.Production
	}

	s := snap.Client{}
	s.New(serverKey, env)
	s.Options.SetPaymentOverrideNotification(baseUrl + "/payment-gateway/midtrans/notification")

	chargeReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  sale.InvoiceID,
			GrossAmt: int64(utils.RoundFloat(sale.TotalPaid)),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: sale.Customer.FirstName,
			LName: sale.Customer.LastName,
			Email: sale.Customer.Email,
			Phone: sale.Customer.Phone,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		Gopay: &snap.GopayDetails{
			EnableCallback: true,
			CallbackUrl:    baseUrl,
		},
		ShopeePay: &snap.ShopeePayDetails{
			CallbackUrl: baseUrl,
		},
		Callbacks: &snap.Callbacks{
			Finish: baseUrl,
		},
		Expiry: &snap.ExpiryDetails{
			StartTime: time.Now().Format("2006-01-02 15:04:05 +0700"),
			Duration:  1,
			Unit:      "day",
		},
		Cstore: &snap.Cstore{
			AlfamartFreeText1: "Invoice ID: " + sale.InvoiceID,
			AlfamartFreeText2: "Generate By: Sales Demo API",
		},
	}

	var totalQuantity int
	for _, item := range sale.Details {
		totalQuantity += item.Quantity
	}

	var items []midtrans.ItemDetails
	for _, data := range sale.Details {
		if sale.Discount > 0 {
			discountAmount := (sale.Discount * sale.Subtotal) / 100
			discountPerItem := discountAmount / float64(totalQuantity)
			if discountPerItem > 0 {
				data.Price -= discountPerItem
			}
		}
		if sale.Tax > 0 {
			taxAmount := (sale.Tax * sale.Subtotal) / 100
			taxPerItem := taxAmount / float64(totalQuantity)
			if taxPerItem > 0 {
				data.Price += taxPerItem
			}
		}
		if sale.MiscPrice > 0 {
			miscPricePerItem := sale.MiscPrice / float64(totalQuantity)
			if miscPricePerItem > 0 {
				data.Price += miscPricePerItem
			}
		}

		item := midtrans.ItemDetails{
			ID:    data.ProductID.String(),
			Name:  data.Product.Name,
			Qty:   int32(data.Quantity),
			Price: int64(utils.RoundFloat(data.Price)),
		}

		items = append(items, item)
	}
	if len(items) > 0 {
		chargeReq.Items = &items
	}

	midtransResponse, midtransError := s.CreateTransaction(chargeReq)
	if midtransError != nil {
		err = errors.New("Failed to request payment to midtrans: " + midtransError.RawError.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		err = errors.New("failed to get location time: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	rawResponse, err := json.Marshal(midtransResponse)
	if err != nil {
		err = errors.New("failed to marshal response: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data := models.SDAMidtransSalePayment{
		SaleID:          sale.ID,
		PaymentMethodID: nil,
		ReferenceCode:   midtransResponse.Token,
		ExpiryDate:      null.TimeFrom(time.Now().In(location).Add(24 * time.Hour)),
		RawResponse:     string(rawResponse),
		UserID:          parsedUserUUID,
		MerchantID:      config.LoadConfig().MidtransMerchantID,
		RedirectUrl:     midtransResponse.RedirectURL,
	}

	response, err = repository.CreateMidtransSalePayment(data)
	if err != nil {
		err = errors.New("failed to create data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusCreated
	return
}

func MidtransHandleNotification(request dto.MidtransNotificationRequest) (statusCode int, err error) {
	expectedSignatureKeyRaw := request.OrderID + request.StatusCode + request.GrossAmount + config.LoadConfig().MidtransServerKey
	hash := sha512.New()
	hash.Write([]byte(expectedSignatureKeyRaw))

	expectedSignatureKey := hex.EncodeToString(hash.Sum(nil))

	if !(expectedSignatureKey == request.SignatureKey) {
		fmt.Println("hmac doesn't match. Expected: ", expectedSignatureKey, ". Get: ", request.SignatureKey)
		err = errors.New("unauthorized")
		statusCode = http.StatusUnauthorized
		return
	}

	if request.FraudStatus != "" {
		if request.FraudStatus != "accept" {
			err = errors.New("transaction indicated as fraud")
			statusCode = http.StatusBadRequest
			return
		}
	}
	if request.PaymentType == "credit_card" {
		if request.TransactionStatus != "capture" && request.TransactionStatus != "settlement" {
			err = errors.New("transaction not yet settled")
			statusCode = http.StatusBadRequest
			return
		}
	} else {
		if request.TransactionStatus != "settlement" {
			err = errors.New("transaction not yet settled")
			statusCode = http.StatusBadRequest
			return
		}
	}

	sale, _, _, _ := repository.GetSales(dto.FindParameter{
		Filter:       "deleted_at IS NULL AND invoice_id = ?",
		FilterValues: []any{request.OrderID},
	}, []string{})
	if len(sale) == 0 {
		err = errors.New("data not found")
		statusCode = http.StatusNotFound
		return
	}
	if sale[0].ID == uuid.Nil {
		err = errors.New("data not found")
		statusCode = http.StatusNotFound
		return
	}

	sale[0].Status = true
	sale[0].PaymentDate = null.TimeFrom(time.Now())

	_, err = repository.UpdateSale(sale[0])
	if err != nil {
		err = errors.New("failed to update data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}
