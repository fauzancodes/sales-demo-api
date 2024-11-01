package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/repository"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func CreateProductCategory(userID string, request dto.ProductCategoryRequest) (response models.SDAProductCategory, statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = errors.New("failed to parse user UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
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
	if err != nil {
		err = errors.New("failed to create data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusCreated
	return
}

func GetProductCategoryByID(id string, preloadFields []string) (data models.SDAProductCategory, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = errors.New("failed to parse UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}
	data, err = repository.GetProductCategoryByID(parsedUUID, preloadFields)
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if err == gorm.ErrRecordNotFound {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	return
}

func GetProductCategories(name, userID string, param utils.PagingRequest, preloadFields []string) (response utils.PagingResponse, data []models.SDAProductCategory, statusCode int, err error) {
	baseFilter := "deleted_at IS NULL"
	if userID != "" {
		baseFilter += " AND user_id = '" + userID + "'"
	}
	filter := baseFilter

	if name != "" {
		filter += " AND name = '" + name + "'"
	}
	if param.Custom != "" {
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

func UpdateProductCategory(id string, request dto.ProductCategoryRequest) (response models.SDAProductCategory, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = errors.New("failed to parse UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data, err := repository.GetProductCategoryByID(parsedUUID, []string{})
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if err == gorm.ErrRecordNotFound {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
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
	if err != nil {
		err = errors.New("failed to update data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func DeleteProductCategory(id string) (statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = errors.New("failed to parse UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data, err := repository.GetProductCategoryByID(parsedUUID, []string{})
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if err == gorm.ErrRecordNotFound {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	err = repository.DeleteProductCategory(data)
	if err != nil {
		err = errors.New("failed to delete data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusNoContent
	return
}

func ImportProductCategory(file *multipart.FileHeader, userID string) (responses []models.SDAProductCategory, statusCode int, err error) {
	rows, err := utils.ValidateImportFile(file, 3)
	if err != nil {
		statusCode = http.StatusBadRequest
		return
	}

	rows = rows[1:]
	if len(rows) == 0 {
		err = errors.New("there is no data in the file")
		statusCode = http.StatusBadRequest
		return
	}

	for _, data := range rows {
		var response models.SDAProductCategory

		check, _, _, _ := repository.GetProductCategories(dto.FindParameter{
			Filter: "deleted_at IS NULL AND code = '" + data[2] + "'",
		}, []string{})

		if len(check) > 0 {
			response = check[0]
		} else {
			check, _, _, _ = repository.GetProductCategories(dto.FindParameter{
				Filter: "deleted_at IS NULL AND name = '" + data[0] + "'",
			}, []string{})

			if len(check) > 0 {
				response = check[0]
			}
		}

		if response.ID == uuid.Nil {
			input := dto.ProductCategoryRequest{
				Code:        data[2],
				Name:        data[0],
				Description: data[1],
				Status:      true,
			}
			if input.Code == "" || input.Code == "-" {
				input.Code = utils.GenerateRandomNumber(12)
			}

			response, statusCode, err = CreateProductCategory(userID, input)
			if err != nil {
				return
			}
		}

		responses = append(responses, response)
	}

	return
}

func ExportProductCategory(userID, fileExtentison string) (filename string, statusCode int, err error) {
	filename = fmt.Sprintf("assets/download/product_categories_%v.%v", userID, fileExtentison)

	productCategories, _, _, err := repository.GetProductCategories(dto.FindParameter{
		BaseFilter: "deleted_at IS NULL AND user_id = '" + userID + "'",
	}, []string{})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	if fileExtentison == "xlsx" {
		f := excelize.NewFile()
		sheetName := "Product Categories"
		f.SetSheetName("Sheet1", sheetName)

		f.SetCellValue(sheetName, "A1", "Name")
		f.SetCellValue(sheetName, "B1", "Description")
		f.SetCellValue(sheetName, "C1", "Code")

		for i, item := range productCategories {
			row := i + 2
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "-")
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "-")
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "-")

			if item.Name != "" {
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), item.Name)
			}
			if item.Description != "" {
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), item.Description)
			}
			if item.Code != "" {
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), item.Code)
			}
		}

		err = f.SaveAs(filename)
		if err != nil {
			statusCode = http.StatusInternalServerError
			err = errors.New("failed to save excel file: " + err.Error())
		}
	}

	if fileExtentison == "csv" {
		var file *os.File
		file, err = os.Create(filename)
		if err != nil {
			statusCode = http.StatusInternalServerError
			err = errors.New("failed to create csv file: " + err.Error())
			return
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		header := []string{"Name", "Description", "Code"}
		err = writer.Write(header)
		if err != nil {
			statusCode = http.StatusInternalServerError
			err = errors.New("failed to write header into csv file: " + err.Error())
			return
		}

		for _, item := range productCategories {
			name := "-"
			description := "-"
			code := "-"

			if item.Name != "" {
				name = item.Name
			}
			if item.Description != "" {
				description = item.Description
			}
			if item.Code != "" {
				code = item.Code
			}

			row := []string{name, description, code}
			err = writer.Write(row)
			if err != nil {
				statusCode = http.StatusInternalServerError
				err = errors.New("failed to write data into csv file: " + err.Error())
				return
			}
		}
	}

	return
}
