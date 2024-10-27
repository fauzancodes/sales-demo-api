package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/pkg/upload"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/repository"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

func GetIPaymuPaymentMethods(code string, param utils.PagingRequest) (response utils.PagingResponse, data []models.SDAIPaymuPaymentMethod, statusCode int, err error) {
	baseFilter := "deleted_at IS NULL"
	filter := baseFilter

	if code != "" {
		filter += " AND code = '" + code + "'"
	}
	if param.Search != "" {
		filter += " AND (name ILIKE '%" + param.Search + "%' OR description ILIKE '%" + param.Search + "%')"
	}

	data, total, totalFiltered, err := repository.GetIPaymuPaymentMethods(dto.FindParameter{
		BaseFilter: baseFilter,
		Filter:     filter,
		Limit:      param.Limit,
		Order:      param.Order,
		Offset:     param.Offset,
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

func GenerateIPaymuSignature(method, ipaymuVA, ipaymuApiKey string, request []byte) (signature string) {
	requestHash := sha256.Sum256(request)
	requestHashToString := hex.EncodeToString(requestHash[:])
	stringToSign := method + ":" + ipaymuVA + ":" + strings.ToLower(requestHashToString) + ":" + ipaymuApiKey

	h := hmac.New(sha256.New, []byte(ipaymuApiKey))
	h.Write([]byte(stringToSign))
	signature = hex.EncodeToString(h.Sum(nil))

	return
}

func SendIPaymuRequest(method, path string, request []byte) (responseBody []byte, err error) {
	var ipaymuVA = config.LoadConfig().IPaymuVA
	var ipaymuApiKey = config.LoadConfig().IPaymuApiKey

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		location = time.Local
	}

	url, err := url.Parse(config.LoadConfig().IPaymuBaseURL + path)
	if err != nil {
		err = errors.New("failed to parse URL: " + err.Error())
		return
	}

	signature := GenerateIPaymuSignature(method, ipaymuVA, ipaymuApiKey, request)

	reqBody := io.NopCloser(strings.NewReader(string(request)))
	currentTime := time.Now().In(location)
	years := currentTime.Year()
	months := int(currentTime.Month())
	days := currentTime.Day()
	hours := currentTime.Hour()
	minutes := currentTime.Minute()
	seconds := currentTime.Second()

	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
			"va":           {ipaymuVA},
			"signature":    {signature},
			"timestamp":    {fmt.Sprintf("%04d%02d%02d%02d%02d%02d", years, months, days, hours, minutes, seconds)},
		},
		Body: reqBody,
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.New("failed to sent request: " + err.Error())
		return
	}
	defer response.Body.Close()

	responseBody, err = io.ReadAll(response.Body)
	if err != nil {
		err = errors.New("failed to read response: " + err.Error())
		return
	}

	return
}

func IPaymuCharge(userID, baseUrl string, request dto.IPaymuSaleRequest) (response models.SDAIPaymuSalePayment, statusCode int, err error) {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		location = time.Local
	}

	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = errors.New("failed to parse user UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	paymentMethodData, _, _, err := repository.GetIPaymuPaymentMethods(dto.FindParameter{
		Filter: "deleted_at IS NULL AND code = '" + strings.ToLower(request.PaymentMethodCode) + "'",
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
		Filter: "deleted_at IS NULL AND invoice_id = '" + request.InvoiceID + "'",
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

	notifyUrl := baseUrl + "/payment-gateway/ipaymu/notification"

	var paymentMethodRequest string
	switch paymentMethod.Code {
	case "bag", "bca", "bni", "cimb", "mandiri", "bmi", "bri", "bsi", "permata", "danamon", "bpd_bali":
		paymentMethodRequest = "va"
	case "alfamart", "indomaret":
		paymentMethodRequest = "cstore"
	case "rpx":
		paymentMethodRequest = "cod"
	case "mpm":
		paymentMethodRequest = "qris"
	case "cc":
		paymentMethodRequest = "cc"
	}

	var totalQuantity int
	for _, item := range sale.Details {
		totalQuantity += item.Quantity
	}

	var productNames []string
	var productQuantities []int
	var productPrices []int64
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

		productNames = append(productNames, data.Product.Name)
		productQuantities = append(productQuantities, data.Quantity)
		productPrices = append(productPrices, int64(utils.RoundFloat(data.Product.Price)))
	}

	regex := regexp.MustCompile("[^0-9]+")
	phoneSanitized := regex.ReplaceAllString(sale.Customer.Phone, "")
	postBodyRaw := dto.IPaymuRequest{
		Name:            sale.Customer.FirstName + " " + sale.Customer.LastName,
		Phone:           phoneSanitized,
		Email:           strings.ToLower(sale.Customer.Email),
		Amount:          int64(utils.RoundFloat(sale.TotalPaid)),
		NotifyURL:       notifyUrl,
		ReferenceID:     sale.InvoiceID,
		PaymentMethod:   paymentMethodRequest,
		PaymentChannel:  paymentMethod.Code,
		ProductName:     productNames,
		ProductQuantity: productQuantities,
		ProductPrice:    productPrices,
		Expired:         24,
	}
	switch paymentMethod.Code {
	case "bsi":
		postBodyRaw.Expired = 3
	case "bri":
		postBodyRaw.Expired = 2
	case "bca":
		postBodyRaw.Expired = 12
	}

	postBody, err := json.Marshal(postBodyRaw)
	if err != nil {
		err = errors.New("failed to marshal post body: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	ipaymuResponse, err := SendIPaymuRequest("POST", "/payment/direct", postBody)
	if err != nil {
		statusCode = http.StatusInternalServerError
		return
	}

	var responseBody dto.IPaymuResponse
	err = json.Unmarshal(ipaymuResponse, &responseBody)
	if err != nil {
		err = errors.New("failed to unmarshal response body: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	if responseBody.Status == 200 {
		data := models.SDAIPaymuSalePayment{
			SaleID:          sale.ID,
			PaymentMethodID: paymentMethod.ID,
			ReferenceCode:   responseBody.Data.TransactionID,
			ExpiryDate:      null.TimeFrom(time.Now().In(location).Add(24 * time.Hour)),
			RawResponse:     string(ipaymuResponse),
			UserID:          parsedUserUUID,
		}

		switch paymentMethod.Code {
		case "bsi":
			data.ExpiryDate = null.TimeFrom(time.Now().In(location).Add(3 * time.Hour))
		case "bri":
			data.ExpiryDate = null.TimeFrom(time.Now().In(location).Add(2 * time.Hour))
		case "bca":
			data.ExpiryDate = null.TimeFrom(time.Now().In(location).Add(12 * time.Hour))
		}

		if paymentMethodRequest == "cc" {
			data.RedirectUrl = responseBody.Data.PaymentCode
		} else if paymentMethodRequest == "qris" {
			var qr *qrcode.QRCode
			qr, err = qrcode.New(responseBody.Data.PaymentCode, qrcode.Medium)
			if err != nil {
				err = errors.New("failed to generate qr code: " + err.Error())
				statusCode = http.StatusInternalServerError
				return
			}

			var buf bytes.Buffer
			err = png.Encode(&buf, qr.Image(256))
			if err != nil {
				err = errors.New("failed to encode qr code to png: " + err.Error())
				statusCode = http.StatusInternalServerError
				return
			}

			data.QRCodeUrl, _, _, err = upload.UploadFile(bytes.NewReader(buf.Bytes()), userID, "")
			if err != nil {
				statusCode = http.StatusInternalServerError
				return
			}
		} else {
			data.PaymentCode = responseBody.Data.PaymentCode
		}

		response, err = repository.CreateIPaymuSalePayment(data)
		if err != nil {
			err = errors.New("failed to create data: " + err.Error())
			statusCode = http.StatusInternalServerError
			return
		}
	} else {
		err = errors.New("failed to charge to ipaymu: " + responseBody.Message)
		statusCode = http.StatusInternalServerError

		return
	}

	statusCode = http.StatusCreated
	return
}

func IPaymuHandleNotification(request dto.IPaymuNotificationRequest) (statusCode int, err error) {
	if request.StatusCode != "1" || request.Status != "berhasil" {
		err = errors.New("transaction not yet settled")
		statusCode = http.StatusBadRequest
		return
	}

	sale, _, _, _ := repository.GetSales(dto.FindParameter{
		Filter: "deleted_at IS NULL AND invoice_id = '" + request.ReferenceID + "'",
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
