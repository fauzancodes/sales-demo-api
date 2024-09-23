package service

import (
	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/repository"
	"github.com/fauzancodes/sales-demo-api/pkg/utils"
	"github.com/google/uuid"
)

func CreateProductCategory(userID string, request dto.ProductCategoryRequest) (response models.SDAProductCategory, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		return
	}

	data := models.SDAProductCategory{
		Name:        request.Name,
		Description: request.Description,
		Status:      request.Status,
		UserID:      parsedUserUUID,
	}

	response, err = repository.CreateProductCategory(data)

	return
}

func GetProductCategoryByID(id string) (data models.SDAProductCategory, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return
	}
	data, err = repository.GetProductCategoryByID(parsedUUID)

	return
}

func GetProductCategories(name string, userID string, param utils.PagingRequest) (response utils.PagingResponse, data []models.SDAProductCategory, err error) {
	baseFilter := "deleted_at IS NULL"
	if userID != "" {
		baseFilter += " AND user_id = '" + userID + "'"
	}
	filter := baseFilter

	if name != "" {
		filter += " AND name = '" + name + "'"
	}
	if param.Custom.(string) != "" {
		filter += " AND status = " + param.Custom.(string)
	}
	if param.Search != "" {
		filter += " AND (name ILIKE '%" + param.Search + "%' OR description ILIKE '%" + param.Search + "%')"
	}

	data, total, totalFiltered, err := repository.GetProductCategories(dto.FindParameter{
		BaseFilter: baseFilter,
		Filter:     filter,
		Limit:      param.Limit,
		Order:      param.Order,
		Offset:     param.Offset,
	})
	if err != nil {
		return
	}

	response = utils.PopulateResPaging(&param, data, total, totalFiltered)

	return
}

func UpdateProductCategory(id string, request dto.ProductCategoryRequest) (response models.SDAProductCategory, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return
	}
	data, err := repository.GetProductCategoryByID(parsedUUID)
	if err != nil {
		return
	}

	if request.Name != "" {
		data.Name = request.Name
	}
	if request.Description != "" {
		data.Description = request.Description
	}
	data.Status = request.Status

	response, err = repository.UpdateProductCategory(data)

	return
}

func DeleteProductCategory(id string) (err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return
	}

	data, err := repository.GetProductCategoryByID(parsedUUID)
	if err != nil {
		return
	}

	err = repository.DeleteProductCategory(data)

	return
}
