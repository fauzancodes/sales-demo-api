package service

import (
	"errors"
	"net/http"
	"strings"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateProductStock(userID string, request dto.ProductStockRequest) (response models.SDAProductStock, statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = errors.New("failed to parse user UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	parsedProductUUID, err := uuid.Parse(request.ProductID)
	if err != nil {
		err = errors.New("failed to parse product UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	lastStock, _ := repository.GetLastProductStock(parsedProductUUID, []string{})

	data := models.SDAProductStock{
		ProductID:   parsedProductUUID,
		Description: request.Description,
		UserID:      parsedUserUUID,
	}

	if strings.ToLower(request.Action) == "add" {
		data.Addition = request.Amount
		data.Current = lastStock.Current + request.Amount
	} else if strings.ToLower(request.Action) == "reduce" {
		data.Reduction = request.Amount
		data.Current = lastStock.Current - request.Amount
	}

	response, err = repository.CreateProductStock(data)
	if err != nil {
		err = errors.New("failed to create data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusCreated
	return
}

func GetProductStocks(productID, userID string, param utils.PagingRequest, preloadFields []string) (response utils.PagingResponse, data []models.SDAProductStock, statusCode int, err error) {
	baseFilter := "deleted_at IS NULL"
	var baseFilterValues []any
	if userID != "" {
		baseFilter += " AND user_id = ?"
		baseFilterValues = append(baseFilterValues, userID)
	}
	filter := baseFilter
	filterValues := baseFilterValues

	if productID != "" {
		filter += " AND product_id = ?"
		filterValues = append(filterValues, productID)
	}

	data, total, totalFiltered, err := repository.GetProductStocks(dto.FindParameter{
		BaseFilter:       baseFilter,
		BaseFilterValues: baseFilterValues,
		Filter:           filter,
		FilterValues:     filterValues,
		Limit:            param.Limit,
		Order:            param.Order,
		Offset:           param.Offset,
	}, preloadFields)
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
