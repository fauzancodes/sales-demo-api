package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/repository"
	"github.com/google/uuid"
	"github.com/guregu/null"
)

func CheckSaleDetails(request dto.SaleRequest) (dto.SaleRequest, error) {
	var err error
	for _, checker := range request.Details {
		var totalDuplicateProduct int
		for _, checked := range request.Details {
			if checker.ProductID == checked.ProductID {
				totalDuplicateProduct++
			}
		}
		if totalDuplicateProduct > 1 {
			err = fmt.Errorf("there are duplicate products. Product ID: %v", checker.ProductID)

			return request, err
		}
	}

	var expectedSubtotal float64
	for _, item := range request.Details {
		expectedSubtotal += float64(item.Quantity) * item.Price

		var parsedProductUUID uuid.UUID
		parsedProductUUID, err = uuid.Parse(item.ProductID)
		if err != nil {
			return request, err
		}

		lastStock, _ := repository.GetLastProductStock(parsedProductUUID, []string{})
		if item.Quantity > lastStock.Current {
			err = fmt.Errorf("insufficient stock. Product ID: %v. Current stock: %v", item.ProductID, lastStock.Current)

			return request, err
		}
	}
	if request.Subtotal == 0 {
		request.Subtotal = expectedSubtotal
	}
	if expectedSubtotal != request.Subtotal {
		err = fmt.Errorf("subtotal does not match. Expected subtotal: %.2f. Just leave it blank if you don't want to bother. It will be calculated automatically", expectedSubtotal)

		return request, err
	}

	totalTax := (request.Tax / 100) * request.Subtotal
	totalDiscount := (request.Discount / 100) * request.Subtotal
	expectedTotalPaid := request.Subtotal + request.MiscPrice + totalTax - totalDiscount
	if request.TotalPaid == 0 {
		request.TotalPaid = expectedTotalPaid
	}
	if expectedTotalPaid != request.TotalPaid {
		err = fmt.Errorf("total_paid does not match. Expected total_paid: %.2f. Just leave it blank if you don't want to bother. It will be calculated automatically. Formula: subtotal + misc_price + (subtotal * (tax / 100)) - (subtotal * (discount / 100))", expectedTotalPaid)

		return request, err
	}

	if request.InvoiceID == "" {
		request.InvoiceID = "INV" + utils.GenerateRandomNumber(12)
	}

	if request.TransactionDate == "" {
		request.TransactionDate = time.Now().Format(time.DateOnly)
	}

	return request, err
}

func CreateSale(userID string, request dto.SaleRequest) (response models.SDASale, err error) {
	request, err = CheckSaleDetails(request)
	if err != nil {
		return
	}

	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		return
	}
	parsedCustomerUUID, err := uuid.Parse(request.CustomerID)
	if err != nil {
		return
	}

	transactionDate, err := time.Parse(time.DateOnly, request.TransactionDate)
	if err != nil {
		return
	}

	data := models.SDASale{
		InvoiceID:       request.InvoiceID,
		Discount:        request.Discount,
		Tax:             request.Tax,
		MiscPrice:       request.MiscPrice,
		Subtotal:        request.Subtotal,
		TotalPaid:       request.TotalPaid,
		TransactionDate: null.TimeFrom(transactionDate),
		UserID:          parsedUserUUID,
		CustomerID:      parsedCustomerUUID,
	}

	response, err = repository.CreateSale(data)
	if err != nil {
		return
	}

	for _, item := range request.Details {
		var parsedProductUUID uuid.UUID
		parsedProductUUID, err = uuid.Parse(item.ProductID)
		if err != nil {
			return
		}

		detail := models.SDASaleDetail{
			ProductID: parsedProductUUID,
			Price:     item.Price,
			Quantity:  item.Quantity,
			SaleID:    response.ID,
			UserID:    parsedUserUUID,
		}

		var detailResponse models.SDASaleDetail
		detailResponse, err = repository.CreateSaleDetail(detail)
		if err != nil {
			return
		}

		lastStock, _ := repository.GetLastProductStock(detailResponse.ProductID, []string{})
		_, err = repository.CreateProductStock(models.SDAProductStock{
			ProductID:   detailResponse.ProductID,
			Reduction:   detailResponse.Quantity,
			Current:     lastStock.Current - detailResponse.Quantity,
			Description: "Stock reduction from create sales " + response.InvoiceID,
			UserID:      response.UserID,
		})
		if err != nil {
			return
		}

		response.Details = append(response.Details, models.SaleDetailRelation{
			CustomGormModel: detailResponse.CustomGormModel,
			ProductID:       detailResponse.ProductID,
			Price:           detailResponse.Price,
			Quantity:        detailResponse.Quantity,
			UserID:          detailResponse.UserID,
			SaleID:          detailResponse.SaleID,
		})
	}

	return
}

func GetSaleByID(id string, preloadFields []string) (response models.SDASale, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return
	}

	response, err = repository.GetSaleByID(parsedUUID, preloadFields)
	if err != nil {
		return
	}

	return
}

func GetSales(invoiceID, userID, customerID, transactionDateMarginTop, transactionDateMarginBottom, productID string, param utils.PagingRequest, preloadFields []string) (response utils.PagingResponse, data []models.SDASale, err error) {
	baseFilter := "deleted_at IS NULL"
	if userID != "" {
		baseFilter += " AND user_id = '" + userID + "'"
	}
	filter := baseFilter

	if invoiceID != "" {
		filter += " AND invoice_id = '" + invoiceID + "'"
	}
	if customerID != "" {
		filter += " AND customer_id = '" + customerID + "'"
	}
	if transactionDateMarginTop != "" {
		filter += " AND transaction_date::DATE <= '" + transactionDateMarginTop + "'"
	}
	if transactionDateMarginBottom != "" {
		filter += " AND transaction_date::DATE >= '" + transactionDateMarginBottom + "'"
	}
	if productID != "" {
		filter += `
			AND id IN(
				SELECT sale_id
				FROM ` + models.SDASaleDetail{}.TableName() + `
				WHERE product_id = '` + productID + `'
			)
		`
	}
	if param.Custom.(string) != "" {
		filter += " AND status = " + param.Custom.(string)
	}
	if param.Search != "" {
		filter += " AND invoice_id ILIKE '%" + param.Search + "%')"
	}

	data, total, totalFiltered, err := repository.GetSales(dto.FindParameter{
		BaseFilter: baseFilter,
		Filter:     filter,
		Limit:      param.Limit,
		Order:      param.Order,
		Offset:     param.Offset,
	}, preloadFields)
	if err != nil {
		return
	}

	response = utils.PopulateResPaging(&param, data, total, totalFiltered)

	return
}

func UpdateSale(id string, request dto.SaleRequest) (response models.SDASale, err error) {
	request, err = CheckSaleDetails(request)
	if err != nil {
		return
	}

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return
	}

	data, err := repository.GetSaleByID(parsedUUID, []string{})
	if err != nil {
		return
	}

	if request.InvoiceID != "" {
		data.InvoiceID = request.InvoiceID
	}
	if request.CustomerID != "" {
		var parsedCustomerUUID uuid.UUID
		parsedCustomerUUID, err = uuid.Parse(request.CustomerID)
		if err != nil {
			return
		}
		data.CustomerID = parsedCustomerUUID
	}
	if request.Discount > 0 {
		data.Discount = request.Discount
	}
	if request.Tax > 0 {
		data.Tax = request.Tax
	}
	if request.MiscPrice > 0 {
		data.MiscPrice = request.MiscPrice
	}
	if request.Subtotal > 0 {
		data.Subtotal = request.Subtotal
	}
	if request.TotalPaid > 0 {
		data.TotalPaid = request.TotalPaid
	}
	if request.TransactionDate != "" {
		var transactionDate time.Time
		transactionDate, err = time.Parse(time.DateOnly, request.TransactionDate)
		if err != nil {
			return
		}
		data.TransactionDate = null.TimeFrom(transactionDate)
	}

	response, err = repository.UpdateSale(data)
	if err != nil {
		return
	}

	if len(request.Details) > 0 {
		var dataDetails []models.SDASaleDetail
		dataDetails, _, _, err = repository.GetSaleDetails(dto.FindParameter{
			BaseFilter: "deleted_at IS NULL AND user_id = '" + response.UserID.String() + "'",
			Filter:     "deleted_at IS NULL AND user_id = '" + data.UserID.String() + "' AND sale_id = '" + data.ID.String() + "'",
		}, []string{})
		if err != nil {
			return
		}

		if len(dataDetails) > 0 {
			for _, item := range dataDetails {
				err = repository.DeleteSaleDetail(item)
				if err != nil {
					return
				}

				lastStock, _ := repository.GetLastProductStock(item.ProductID, []string{})
				_, err = repository.CreateProductStock(models.SDAProductStock{
					ProductID:   item.ProductID,
					Addition:    item.Quantity,
					Current:     lastStock.Current + item.Quantity,
					Description: "Stock addition from update sales " + data.InvoiceID,
					UserID:      data.UserID,
				})
				if err != nil {
					return
				}
			}
		}

		for _, item := range request.Details {
			var parsedProductUUID uuid.UUID
			parsedProductUUID, err = uuid.Parse(item.ProductID)
			if err != nil {
				return
			}

			detail := models.SDASaleDetail{
				ProductID: parsedProductUUID,
				Price:     item.Price,
				Quantity:  item.Quantity,
				SaleID:    response.ID,
				UserID:    response.UserID,
			}

			var detailResponse models.SDASaleDetail
			detailResponse, err = repository.CreateSaleDetail(detail)
			if err != nil {
				return
			}

			lastStock, _ := repository.GetLastProductStock(detailResponse.ProductID, []string{})
			_, err = repository.CreateProductStock(models.SDAProductStock{
				ProductID:   detailResponse.ProductID,
				Reduction:   detailResponse.Quantity,
				Current:     lastStock.Current - detailResponse.Quantity,
				Description: "Stock reduction from update sales " + response.InvoiceID,
				UserID:      response.UserID,
			})
			if err != nil {
				return
			}

			response.Details = append(response.Details, models.SaleDetailRelation{
				CustomGormModel: detailResponse.CustomGormModel,
				ProductID:       detailResponse.ProductID,
				Price:           detailResponse.Price,
				Quantity:        detailResponse.Quantity,
				UserID:          detailResponse.UserID,
				SaleID:          detailResponse.SaleID,
			})
		}
	}

	return
}

func DeleteSale(id string) (err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return
	}

	data, err := repository.GetSaleByID(parsedUUID, []string{})
	if err != nil {
		return
	}

	dataDetails, _, _, err := repository.GetSaleDetails(dto.FindParameter{
		BaseFilter: "deleted_at IS NULL AND user_id = '" + data.UserID.String() + "'",
		Filter:     "deleted_at IS NULL AND user_id = '" + data.UserID.String() + "' AND sale_id = '" + data.ID.String() + "'",
	}, []string{})
	if err != nil {
		return
	}

	if len(dataDetails) > 0 {
		for _, item := range dataDetails {
			err = repository.DeleteSaleDetail(item)
			if err != nil {
				return
			}

			lastStock, _ := repository.GetLastProductStock(item.ProductID, []string{})
			_, err = repository.CreateProductStock(models.SDAProductStock{
				ProductID:   item.ProductID,
				Addition:    item.Quantity,
				Current:     lastStock.Current + item.Quantity,
				Description: "Stock addition from delete sales " + data.InvoiceID,
				UserID:      data.UserID,
			})
			if err != nil {
				return
			}
		}
	}

	err = repository.DeleteSale(data)

	return
}

func SendSaleInvoice(saleID uuid.UUID) {
	sale, _ := repository.GetSaleByID(saleID, []string{"Details"})
	if sale.ID != uuid.Nil {
		fill := dto.SaleInvoice{
			InvoiceID:       sale.InvoiceID,
			TransactionDate: sale.TransactionDate.Time.Format(time.RFC1123Z),
			Status:          "Unpaid",
			Subtotal:        sale.Subtotal,
			Discount:        sale.Discount,
			Tax:             sale.Tax,
			MiscPrice:       sale.MiscPrice,
			TotalPaid:       sale.TotalPaid,
		}

		var customer models.SDACustomer
		if sale.CustomerID != uuid.Nil {
			customer, _ = repository.GetCustomerByID(sale.CustomerID, []string{})
		}
		if customer.ID != uuid.Nil {
			fill.CustomerFullname = customer.FirstName + " " + customer.LastName
		}
		if strings.ReplaceAll(fill.CustomerFullname, " ", "") == "" {
			fill.CustomerFullname = customer.Email
		}

		var user models.SDAUser
		if sale.UserID != uuid.Nil {
			user, _ = repository.GetUserByID(sale.UserID, []string{})
		}
		if user.ID != uuid.Nil {
			fill.UserFullname = user.FirstName + " " + user.LastName
		}
		if strings.ReplaceAll(fill.UserFullname, " ", "") == "" {
			fill.UserFullname = user.Email
		}

		if sale.Status {
			fill.Status = "Paid"
		}

		if len(sale.Details) > 0 {
			for _, detail := range sale.Details {
				fillDetail := dto.SaleInvoiceDetail{
					Quantity:     detail.Quantity,
					ProductPrice: detail.Price,
				}
				fillDetail.TotalPrice = float64(fillDetail.Quantity) * fillDetail.ProductPrice

				var product models.SDAProduct
				if detail.ProductID != uuid.Nil {
					product, _ = repository.GetProductByID(detail.ProductID, []string{})
				}
				if product.ID != uuid.Nil {
					fillDetail.ProductName = product.Name
				}

				fill.Details = append(fill.Details, fillDetail)
			}
		}

		utils.SendEmail("invoice", user.Email, customer.Email, "Sales Invoice", "", fill)
	}
}
