package service

import (
	"encoding/csv"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strconv"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/repository"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

func CreateProductCategory(userID string, request dto.ProductCategoryRequest) (response models.SDAProductCategory, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		return
	}
	if request.Code == "" || request.Code == "-" {
		request.Code = utils.GenerateRandomNumber(12)
	}

	data := models.SDAProductCategory{
		Code:        request.Code,
		Name:        request.Name,
		Description: request.Description,
		Status:      request.Status,
		UserID:      parsedUserUUID,
	}

	response, err = repository.CreateProductCategory(data)

	return
}

func GetProductCategoryByID(id string, preloadFields []string) (data models.SDAProductCategory, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return
	}
	data, err = repository.GetProductCategoryByID(parsedUUID, preloadFields)

	return
}

func GetProductCategories(name, userID string, param utils.PagingRequest, preloadFields []string) (response utils.PagingResponse, data []models.SDAProductCategory, err error) {
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
	}, preloadFields)
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

	data, err := repository.GetProductCategoryByID(parsedUUID, []string{})
	if err != nil {
		return
	}

	if request.Code != "" {
		data.Code = request.Code
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

	data, err := repository.GetProductCategoryByID(parsedUUID, []string{})
	if err != nil {
		return
	}

	err = repository.DeleteProductCategory(data)

	return
}

func ImportProductCategory(file *multipart.FileHeader, userID string) (responses []models.SDAProductCategory, err error) {
	src, err := file.Open()
	if err != nil {
		return
	}
	defer src.Close()

	var rows [][]string
	extension := filepath.Ext(file.Filename)

	if extension == ".xls" || extension == ".xlsx" {
		var f *excelize.File
		f, err = excelize.OpenReader(src)
		if err != nil {
			return
		}

		sheets := f.GetSheetList()
		if len(sheets) == 0 {
			err = errors.New("there is no sheet in the file")
			return
		}

		rows, err = f.GetRows(sheets[0])
		if err != nil {
			return
		}
		if len(rows[0]) != 3 {
			err = errors.New("The number of columns must match the template. Expected: 3 columns. Current: " + strconv.Itoa(len(rows[0])) + " columns")
			return
		}
	} else if extension == ".csv" {
		reader := csv.NewReader(src)
		rows, err = reader.ReadAll()
		if err != nil {
			return
		}
	} else {
		err = errors.New("the file format only accepts .xls, .xlsx, .csv")
		return
	}

	rows = rows[1:]
	if len(rows) > 0 {
		for _, data := range rows {
			var response models.SDAProductCategory
			input := dto.ProductCategoryRequest{
				Code:        data[2],
				Name:        data[0],
				Description: data[1],
				Status:      true,
			}
			if input.Code == "" || input.Code == "-" {
				input.Code = utils.GenerateRandomNumber(12)
			}

			response, err = CreateProductCategory(userID, input)
			if err != nil {
				return
			}

			responses = append(responses, response)
		}
	} else {
		err = errors.New("there is no data in the file")
		return
	}

	return
}
