package service

import (
	"strings"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/repository"
	"github.com/google/uuid"
)

func CreateProductStock(userID string, request dto.ProductStockRequest) (response models.SDAProductStock, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		return
	}

	parsedProductUUID, err := uuid.Parse(request.ProductID)
	if err != nil {
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

	return
}

func GetProductStocks(productID, userID string, param utils.PagingRequest, preloadFields []string) (response utils.PagingResponse, data []models.SDAProductStock, err error) {
	baseFilter := "deleted_at IS NULL"
	if userID != "" {
		baseFilter += " AND user_id = '" + userID + "'"
	}
	filter := baseFilter

	if productID != "" {
		filter += " AND product_id = '" + productID + "'"
	}

	data, total, totalFiltered, err := repository.GetProductStocks(dto.FindParameter{
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
