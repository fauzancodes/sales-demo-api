package service

import (
	"encoding/json"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/repository"
	"github.com/fauzancodes/sales-demo-api/pkg/utils"
	"github.com/google/uuid"
)

func BuildProductResponse(data models.SDAProduct) (response dto.ProductResponse, err error) {
	response.CustomGormModel = data.CustomGormModel
	response.Name = data.Name
	response.Description = data.Description
	response.Status = data.Status
	response.Price = data.Price

	err = json.Unmarshal([]byte(data.Image), &response.Image)
	if err != nil {
		return
	}

	category, err := repository.GetProductCategoryByID(data.CategoryID)
	if err != nil {
		return
	}
	response.Category.ID = category.ID
	response.Category.Name = category.Name

	return
}

func CreateProduct(userID string, request dto.ProductRequest) (response models.SDAProduct, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		return
	}

	parsedCategoryUUID, err := uuid.Parse(request.CategoryID)
	if err != nil {
		return
	}

	jsonImage, err := json.Marshal(request.Image)
	if err != nil {
		return
	}

	data := models.SDAProduct{
		Name:        request.Name,
		Description: request.Description,
		Status:      request.Status,
		UserID:      parsedUserUUID,
		CategoryID:  parsedCategoryUUID,
		Price:       request.Price,
		Image:       string(jsonImage),
	}

	response, err = repository.CreateProduct(data)

	return
}

func GetProductByID(id string) (response dto.ProductResponse, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return
	}

	data, err := repository.GetProductByID(parsedUUID)
	if err != nil {
		return
	}

	response, err = BuildProductResponse(data)
	if err != nil {
		return
	}

	return
}

func GetProducts(name, userID, categoryID string, param utils.PagingRequest) (response utils.PagingResponse, data []models.SDAProduct, err error) {
	baseFilter := "deleted_at IS NULL"
	if userID != "" {
		baseFilter += " AND user_id = '" + userID + "'"
	}
	filter := baseFilter

	if name != "" {
		filter += " AND name = '" + name + "'"
	}
	if categoryID != "" {
		filter += " AND category_id = '" + categoryID + "'"
	}
	if param.Custom.(string) != "" {
		filter += " AND status = " + param.Custom.(string)
	}
	if param.Search != "" {
		filter += " AND (name ILIKE '%" + param.Search + "%' OR description ILIKE '%" + param.Search + "%')"
	}

	data, total, totalFiltered, err := repository.GetProducts(dto.FindParameter{
		BaseFilter: baseFilter,
		Filter:     filter,
		Limit:      param.Limit,
		Order:      param.Order,
		Offset:     param.Offset,
	})
	if err != nil {
		return
	}

	var results []dto.ProductResponse
	var result dto.ProductResponse
	for _, item := range data {
		result, err = BuildProductResponse(item)
		if err != nil {
			return
		}

		results = append(results, result)
	}

	response = utils.PopulateResPaging(&param, results, total, totalFiltered)

	return
}

func UpdateProduct(id string, request dto.ProductRequest) (response models.SDAProduct, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return
	}

	data, err := repository.GetProductByID(parsedUUID)
	if err != nil {
		return
	}

	if request.Name != "" {
		data.Name = request.Name
	}
	if request.Description != "" {
		data.Description = request.Description
	}
	if request.Price > 0 {
		data.Price = request.Price
	}

	var parsedCategoryUUID uuid.UUID
	if request.CategoryID != "" {
		parsedCategoryUUID, err = uuid.Parse(request.CategoryID)
		if err != nil {
			return
		}
		data.CategoryID = parsedCategoryUUID
	}

	data.Status = request.Status

	var jsonImage []byte
	if len(request.Image) > 0 {
		jsonImage, err = json.Marshal(request.Image)
		if err != nil {
			return
		}

		data.Image = string(jsonImage)
	}

	response, err = repository.UpdateProduct(data)

	return
}

func DeleteProduct(id string) (err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return
	}

	data, err := repository.GetProductByID(parsedUUID)
	if err != nil {
		return
	}

	err = repository.DeleteProduct(data)

	return
}
