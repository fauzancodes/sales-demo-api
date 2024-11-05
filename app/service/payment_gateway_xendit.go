package service

import (
	"context"
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
	"github.com/xendit/xendit-go/v6"
	"github.com/xendit/xendit-go/v6/invoice"
	"github.com/xendit/xendit-go/v6/payment_request"
	"gorm.io/gorm"
)

func GetXenditPaymentMethods(code string, param utils.PagingRequest) (response utils.PagingResponse, data []models.SDAXenditPaymentMethod, statusCode int, err error) {
	baseFilter := "deleted_at IS NULL"
	filter := baseFilter
	var fileterValues []any

	if code != "" {
		filter += " AND code = ?"
		fileterValues = append(fileterValues, code)
	}
	if param.Search != "" {
		filter += " AND (name ILIKE ? OR description ILIKE ?)"
		fileterValues = append(fileterValues, fmt.Sprintf("%%%s%%", param.Search))
	}

	data, total, totalFiltered, err := repository.GetXenditPaymentMethods(dto.FindParameter{
		BaseFilter:   baseFilter,
		Filter:       filter,
		FilterValues: fileterValues,
		Limit:        param.Limit,
		Order:        param.Order,
		Offset:       param.Offset,
	})
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if err == gorm.ErrRecordNotFound {
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

func XenditChargePayment(userID, baseUrl string, request dto.XenditRequestPayment) (response models.SDAXenditSalePayment, statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = errors.New("failed to parse user UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	paymentMethodData, _, _, err := repository.GetXenditPaymentMethods(dto.FindParameter{
		Filter:       "deleted_at IS NULL AND code = ?",
		FilterValues: []any{strings.ToUpper(request.PaymentMethodCode)},
	})
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if err == gorm.ErrRecordNotFound {
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
		if err == gorm.ErrRecordNotFound {
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

	client := xendit.NewClient(config.LoadConfig().XenditSecretKey)

	amount := float64(utils.RoundFloat(sale.TotalPaid))
	paymentRequestParameters := payment_request.PaymentRequestParameters{
		ReferenceId:   &sale.InvoiceID,
		Amount:        &amount,
		Currency:      payment_request.PAYMENTREQUESTCURRENCY_IDR,
		CaptureMethod: *payment_request.NewNullablePaymentRequestCaptureMethod(payment_request.PAYMENTREQUESTCAPTUREMETHOD_AUTOMATIC.Ptr()),
	}

	var customerPhone string
	if sale.Customer.Phone[:2] == "08" {
		customerPhone = "+62" + sale.Customer.Phone[1:]
	} else if sale.Customer.Phone[:3] == "+62" {
		customerPhone = sale.Customer.Phone
	} else {
		customerPhone = "+62" + sale.Customer.Phone
	}
	customer, _, _ := client.CustomerApi.GetCustomerByReferenceID(context.Background()).ReferenceId(sale.Customer.ID.String()).Execute()
	if len(customer.Data) > 0 {
		paymentRequestParameters.CustomerId = *payment_request.NewNullableString(&customer.Data[0].Id)
	} else {
		paymentRequestParameters.Customer = map[string]interface{}{
			"reference_id": sale.Customer.ID,
			"type":         "INDIVIDUAL",
			"individual_detail": map[string]interface{}{
				"given_names": sale.Customer.FirstName,
				"surname":     sale.Customer.LastName,
			},
			"email":         sale.Customer.Email,
			"mobile_number": customerPhone,
		}
	}

	var totalQuantity int
	for _, item := range sale.Details {
		totalQuantity += item.Quantity
	}

	var paymentRequestItems []payment_request.PaymentRequestBasketItem
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

		var productCategory models.SDAProductCategory
		productCategory, err = repository.GetProductCategoryByID(data.Product.CategoryID, []string{})
		if err != nil {
			err = errors.New("failed to get data: " + err.Error())
			if err == gorm.ErrRecordNotFound {
				statusCode = http.StatusNotFound
				return
			}

			statusCode = http.StatusInternalServerError
			return
		}

		productID := data.ProductID.String()
		productType := "PHYSICAL_PRODUCT"
		paymentRequestItems = append(paymentRequestItems, payment_request.PaymentRequestBasketItem{
			ReferenceId: &productID,
			Type:        &productType,
			Name:        data.Product.Name,
			Category:    productCategory.Name,
			Currency:    payment_request.PAYMENTREQUESTCURRENCY_IDR.String(),
			Quantity:    float64(data.Quantity),
			Price:       float64(utils.RoundFloat(data.Price)),
		})
	}
	paymentRequestParameters.Items = paymentRequestItems

	paymentMethodParameters := payment_request.PaymentMethodParameters{
		Reusability: payment_request.PAYMENTMETHODREUSABILITY_ONE_TIME_USE,
	}

	switch paymentMethod.Code {
	case "DANA", "OVO", "LINKAJA", "ASTRAPAY", "JENIUSPAY", "SHOPEEPAY":
		paymentMethodParameters.Type = payment_request.PAYMENTMETHODTYPE_EWALLET

		var channelCode payment_request.EWalletChannelCode
		channelProperties := payment_request.EWalletChannelProperties{
			SuccessReturnUrl: &request.SuccessReturnUrl,
			FailureReturnUrl: &request.FailedReturnUrl,
		}
		switch paymentMethod.Code {
		case "DANA":
			channelCode = payment_request.EWALLETCHANNELCODE_DANA
		case "OVO":
			channelCode = payment_request.EWALLETCHANNELCODE_OVO
			channelProperties.MobileNumber = &request.EWallet.MobileNumber
		case "LINKAJA":
			channelCode = payment_request.EWALLETCHANNELCODE_LINKAJA
		case "ASTRAPAY":
			channelCode = payment_request.EWALLETCHANNELCODE_ASTRAPAY
		case "JENIUSPAY":
			channelCode = payment_request.EWALLETCHANNELCODE_JENIUSPAY
			channelProperties.Cashtag = &request.EWallet.CashTag
		case "SHOPEEPAY":
			channelCode = payment_request.EWALLETCHANNELCODE_SHOPEEPAY
		}

		paymentMethodParameters.Ewallet = *payment_request.NewNullableEWalletParameters(&payment_request.EWalletParameters{
			ChannelCode:       &channelCode,
			ChannelProperties: &channelProperties,
		})
	case "DIRECT_DEBIT_BRI", "DIRECT_DEBIT_MANDIRI":
		paymentMethodParameters.Type = payment_request.PAYMENTMETHODTYPE_DIRECT_DEBIT

		var channelCode payment_request.DirectDebitChannelCode
		if paymentMethod.Code == "DIRECT_DEBIT_BRI" {
			channelCode = payment_request.DIRECTDEBITCHANNELCODE_BRI
		}
		if paymentMethod.Code == "DIRECT_DEBIT_MANDIRI" {
			channelCode = payment_request.DIRECTDEBITCHANNELCODE_MANDIRI
		}

		directDebitParameterChannelProperties := payment_request.DirectDebitChannelProperties{
			DirectDebitChannelPropertiesBankRedirect: &payment_request.DirectDebitChannelPropertiesBankRedirect{
				Email:            &request.DirectDebit.Email,
				MobileNumber:     &request.DirectDebit.MobileNumber,
				SuccessReturnUrl: &request.SuccessReturnUrl,
				FailureReturnUrl: &request.FailedReturnUrl,
			},
		}

		var debitType payment_request.DirectDebitType
		if request.DirectDebit.Type == "BANK_ACCOUNT" {
			debitType = payment_request.DIRECTDEBITTYPE_BANK_ACCOUNT
			directDebitParameterChannelProperties.DirectDebitChannelPropertiesBankAccount = &payment_request.DirectDebitChannelPropertiesBankAccount{
				SuccessReturnUrl: &request.SuccessReturnUrl,
				FailureReturnUrl: &request.FailedReturnUrl,
				MobileNumber:     &request.DirectDebit.MobileNumber,
			}
		}
		if request.DirectDebit.Type == "DEBIT_CARD" {
			debitType = payment_request.DIRECTDEBITTYPE_DEBIT_CARD
			directDebitParameterChannelProperties.DirectDebitChannelPropertiesDebitCard = &payment_request.DirectDebitChannelPropertiesDebitCard{
				MobileNumber: &request.DirectDebit.MobileNumber,
				CardLastFour: &request.DirectDebit.LastFourDigit,
				CardExpiry:   &request.DirectDebit.CardExpiry,
				Email:        &request.DirectDebit.Email,
			}
		}

		paymentMethodParameters.DirectDebit = *payment_request.NewNullableDirectDebitParameters(&payment_request.DirectDebitParameters{
			Type:              &debitType,
			ChannelCode:       channelCode,
			ChannelProperties: *payment_request.NewNullableDirectDebitChannelProperties(&directDebitParameterChannelProperties),
		})
	case "CARD":
		paymentMethodParameters.Type = payment_request.PAYMENTMETHODTYPE_CARD
		var skip3DSecure bool
		maskedCardNumber := request.Card.CardNumber[:6] + request.Card.CardNumber[len(request.Card.CardNumber)-4:]
		expiryMonth := request.Card.ExpiryMonth
		expiryYear := request.Card.ExpiryYear
		cardHolderName := payment_request.NewNullableString(&request.Card.CardHolderName)
		country := "ID"
		cardNumber := request.Card.CardNumber

		paymentMethodParameters.Card = *payment_request.NewNullableCardParameters(&payment_request.CardParameters{
			ChannelProperties: payment_request.CardChannelProperties{
				SkipThreeDSecure: *payment_request.NewNullableBool(&skip3DSecure),
				SuccessReturnUrl: *payment_request.NewNullableString(&request.SuccessReturnUrl),
				FailureReturnUrl: *payment_request.NewNullableString(&request.FailedReturnUrl),
			},
			CardInformation: &payment_request.CardInformation{
				MaskedCardNumber: &maskedCardNumber,
				ExpiryMonth:      &expiryMonth,
				ExpiryYear:       &expiryYear,
				CardholderName:   *cardHolderName,
				Country:          &country,
				CardNumber:       &cardNumber,
			},
		})
	case "ALFAMART", "INDOMARET":
		paymentMethodParameters.Type = payment_request.PAYMENTMETHODTYPE_OVER_THE_COUNTER

		currency := payment_request.PAYMENTREQUESTCURRENCY_IDR
		var channelCode payment_request.OverTheCounterChannelCode
		if paymentMethod.Code == "ALFAMART" {
			channelCode = payment_request.OVERTHECOUNTERCHANNELCODE_ALFAMART
		}
		if paymentMethod.Code == "INDOMARET" {
			channelCode = payment_request.OVERTHECOUNTERCHANNELCODE_INDOMARET
		}
		expiry := time.Now().Add(24 * time.Hour)
		paymentMethodParameters.OverTheCounter = *payment_request.NewNullableOverTheCounterParameters(&payment_request.OverTheCounterParameters{
			Amount:      *payment_request.NewNullableFloat64(&amount),
			Currency:    &currency,
			ChannelCode: channelCode,
			ChannelProperties: payment_request.OverTheCounterChannelProperties{
				CustomerName: sale.Customer.FirstName + " " + sale.Customer.LastName,
				ExpiresAt:    &expiry,
			},
		})
	case "BCA", "BSI", "BJB", "CIMB", "SAHABAT_SAMPOERNA", "ARTAJASA", "BRI", "BNI", "MANDIRI", "PERMATA":
		paymentMethodParameters.Type = payment_request.PAYMENTMETHODTYPE_VIRTUAL_ACCOUNT

		var channelCode payment_request.VirtualAccountChannelCode
		switch paymentMethod.Code {
		case "BCA":
			channelCode = payment_request.VIRTUALACCOUNTCHANNELCODE_BCA
		case "BSI":
			channelCode = payment_request.VIRTUALACCOUNTCHANNELCODE_BSI
		case "BJB":
			channelCode = payment_request.VIRTUALACCOUNTCHANNELCODE_BJB
		case "CIMB":
			channelCode = payment_request.VIRTUALACCOUNTCHANNELCODE_CIMB
		case "SAHABAT_SAMPOERNA":
			channelCode = payment_request.VIRTUALACCOUNTCHANNELCODE_SAHABAT_SAMPOERNA
		case "ARTAJASA":
			channelCode = payment_request.VIRTUALACCOUNTCHANNELCODE_ARTAJASA
		case "BRI":
			channelCode = payment_request.VIRTUALACCOUNTCHANNELCODE_BRI
		case "BNI":
			channelCode = payment_request.VIRTUALACCOUNTCHANNELCODE_BNI
		case "MANDIRI":
			channelCode = payment_request.VIRTUALACCOUNTCHANNELCODE_MANDIRI
		case "PERMATA":
			channelCode = payment_request.VIRTUALACCOUNTCHANNELCODE_PERMATA
		}

		currency := payment_request.PAYMENTREQUESTCURRENCY_IDR
		expiry := time.Now().Add(24 * time.Hour)
		paymentMethodParameters.VirtualAccount = *payment_request.NewNullableVirtualAccountParameters(&payment_request.VirtualAccountParameters{
			MinAmount:   *payment_request.NewNullableFloat64(&amount),
			MaxAmount:   *payment_request.NewNullableFloat64(&amount),
			Amount:      *payment_request.NewNullableFloat64(&amount),
			Currency:    &currency,
			ChannelCode: channelCode,
			ChannelProperties: payment_request.VirtualAccountChannelProperties{
				CustomerName: sale.Customer.FirstName + " " + sale.Customer.LastName,
				ExpiresAt:    &expiry,
			},
		})
	case "QR_CODE_DANA", "QR_CODE_LINKAJA":
		paymentMethodParameters.Type = payment_request.PAYMENTMETHODTYPE_QR_CODE

		var channelCode payment_request.QRCodeChannelCode
		if paymentMethod.Code == "QR_CODE_DANA" {
			channelCode = payment_request.QRCODECHANNELCODE_DANA
		}
		if paymentMethod.Code == "QR_CODE_LINKAJA" {
			channelCode = payment_request.QRCODECHANNELCODE_LINKAJA
		}

		expiry := time.Now().Add(24 * time.Hour)
		paymentMethodParameters.QrCode = *payment_request.NewNullableQRCodeParameters(&payment_request.QRCodeParameters{
			ChannelCode: *payment_request.NewNullableQRCodeChannelCode(&channelCode),
			ChannelProperties: &payment_request.QRCodeChannelProperties{
				ExpiresAt: &expiry,
			},
		})
	}
	paymentRequestParameters.PaymentMethod = &paymentMethodParameters

	xenditResponse, _, xenditError := client.PaymentRequestApi.CreatePaymentRequest(context.Background()).
		PaymentRequestParameters(paymentRequestParameters).
		Execute()
	if xenditError != nil {
		xenditErrorRawResponse, _ := json.Marshal(xenditError.RawResponse())
		err = errors.New("Failed to request payment to xendit: " + string(xenditErrorRawResponse))
		statusCode = http.StatusInternalServerError
		return
	}
	fmt.Println("xenditResponse:", xenditResponse)
	for i, item := range xenditResponse.Actions {
		if i > 0 {
			if item.UrlType == "" {
				xenditResponse.Actions[i].UrlType = xenditResponse.Actions[i-1].UrlType
			}
		}
	}

	rawResponse, err := json.Marshal(xenditResponse)
	if err != nil {
		err = errors.New("failed to marshall response: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}
	data := models.SDAXenditSalePayment{
		SaleID:          sale.ID,
		PaymentMethodID: &paymentMethod.ID,
		ReferenceCode:   xenditResponse.Id,
		ExpiryDate:      null.TimeFrom(time.Now().Add(24 * time.Hour)),
		RawResponse:     string(rawResponse),
		UserID:          parsedUserUUID,
	}

	if xenditResponse.Status == "REQUIRES_ACTION" {
		for _, item := range xenditResponse.Actions {
			if *item.Url.Get() != "" {
				data.RedirectUrl = *item.Url.Get()
				break
			}
		}
	}
	switch xenditResponse.PaymentMethod.Type {
	case payment_request.PAYMENTMETHODTYPE_CARD:
		//do something later . . . . .
	case payment_request.PAYMENTMETHODTYPE_EWALLET:
		//do something later . . . . .
	case payment_request.PAYMENTMETHODTYPE_DIRECT_DEBIT:
		//do something later . . . . .
	case payment_request.PAYMENTMETHODTYPE_OVER_THE_COUNTER:
		data.PaymentCode = *xenditResponse.PaymentMethod.OverTheCounter.Get().ChannelProperties.PaymentCode
	case payment_request.PAYMENTMETHODTYPE_QR_CODE:
		data.QRCodeUrl = *xenditResponse.PaymentMethod.QrCode.Get().ChannelProperties.QrString
	case payment_request.PAYMENTMETHODTYPE_VIRTUAL_ACCOUNT:
		data.PaymentCode = *xenditResponse.PaymentMethod.VirtualAccount.Get().ChannelProperties.VirtualAccountNumber
	}

	response, err = repository.CreateXenditSalePayment(data)
	if err != nil {
		err = errors.New("failed to create data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusCreated
	return
}

func XenditChargeInvoice(userID, baseUrl string, request dto.XenditRequestInvoice) (response models.SDAXenditSalePayment, statusCode int, err error) {
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
		if err == gorm.ErrRecordNotFound {
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

	client := xendit.NewClient(config.LoadConfig().XenditSecretKey)

	amount := float64(utils.RoundFloat(sale.TotalPaid))

	var customerPhone string
	if sale.Customer.Phone[:2] == "08" {
		customerPhone = "+62" + sale.Customer.Phone[1:]
	} else if sale.Customer.Phone[:3] == "+62" {
		customerPhone = sale.Customer.Phone
	} else {
		customerPhone = "+62" + sale.Customer.Phone
	}

	var totalQuantity int
	for _, item := range sale.Details {
		totalQuantity += item.Quantity
	}

	var paymentRequestItems []invoice.InvoiceItem
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

		var productCategory models.SDAProductCategory
		productCategory, err = repository.GetProductCategoryByID(data.Product.CategoryID, []string{})
		if err != nil {
			err = errors.New("failed to get data: " + err.Error())
			if err == gorm.ErrRecordNotFound {
				statusCode = http.StatusNotFound
				return
			}

			statusCode = http.StatusInternalServerError
			return
		}

		productID := data.ProductID.String()
		paymentRequestItems = append(paymentRequestItems, invoice.InvoiceItem{
			ReferenceId: &productID,
			Name:        data.Product.Name,
			Category:    &productCategory.Name,
			Quantity:    float32(data.Quantity),
			Price:       float32(utils.RoundFloat(data.Price)),
		})
	}

	currency := invoice.INVOICECURRENCY_IDR.String()
	shouldAuthenticateCreditCard := true
	xenditResponse, _, xenditError := client.InvoiceApi.CreateInvoice(context.Background()).
		CreateInvoiceRequest(invoice.CreateInvoiceRequest{
			ExternalId: sale.InvoiceID,
			Amount:     amount,
			Customer: &invoice.CustomerObject{
				PhoneNumber:  *invoice.NewNullableString(&customerPhone),
				GivenNames:   *invoice.NewNullableString(&sale.Customer.FirstName),
				Surname:      *invoice.NewNullableString(&sale.Customer.LastName),
				Email:        *invoice.NewNullableString(&sale.Customer.Email),
				MobileNumber: *invoice.NewNullableString(&customerPhone),
			},
			SuccessRedirectUrl:           &baseUrl,
			Currency:                     &currency,
			Items:                        paymentRequestItems,
			ShouldAuthenticateCreditCard: &shouldAuthenticateCreditCard,
		}).
		Execute()
	if xenditError != nil {
		xenditErrorRawResponse, _ := json.Marshal(xenditError.RawResponse())
		err = errors.New("Failed to request payment to xendit: " + string(xenditErrorRawResponse))
		statusCode = http.StatusInternalServerError
		return
	}

	rawResponse, err := json.Marshal(xenditResponse)
	if err != nil {
		err = errors.New("failed to marshall response: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}
	data := models.SDAXenditSalePayment{
		SaleID:          sale.ID,
		PaymentMethodID: nil,
		ReferenceCode:   *xenditResponse.Id,
		ExpiryDate:      null.TimeFrom(time.Now().Add(24 * time.Hour)),
		RawResponse:     string(rawResponse),
		UserID:          parsedUserUUID,
		RedirectUrl:     xenditResponse.InvoiceUrl,
	}

	response, err = repository.CreateXenditSalePayment(data)
	if err != nil {
		err = errors.New("failed to create data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusCreated
	return
}

func XenditHandleNotification(request dto.XenditNotificationRequest, callbackToken string) (statusCode int, err error) {
	if request.BusinessID != config.LoadConfig().XenditBusinessID {
		fmt.Println("business_id doesn't match: ", request.BusinessID)
		err = errors.New("unauthorized")
		statusCode = http.StatusUnauthorized
		return
	}

	if callbackToken != config.LoadConfig().XenditWebhookToken {
		fmt.Println("x-callback-token doesn't match: ", callbackToken)
		err = errors.New("unauthorized")
		statusCode = http.StatusUnauthorized
		return
	}

	if request.Event != "payment.succeeded" {
		err = errors.New("transaction not yet settled")
		statusCode = http.StatusBadRequest
		return
	}

	sale, _, _, _ := repository.GetSales(dto.FindParameter{
		Filter:       "deleted_at IS NULL AND invoice_id = ?",
		FilterValues: []any{request.Data.ReferenceID},
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
