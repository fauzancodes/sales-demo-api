package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/pkg/upload"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/repository"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func BuildProductResponse(data models.SDAProduct) (response dto.ProductResponse, err error) {
	response.CustomGormModel = data.CustomGormModel
	response.Code = data.Code
	response.Name = data.Name
	response.Description = data.Description
	response.Status = data.Status
	response.Price = data.Price
	if data.Category != nil {
		response.Category = *data.Category
	}

	err = json.Unmarshal([]byte(data.Image), &response.Image)
	if err != nil {
		err = errors.New("failed to unmarshal image: " + err.Error())
		return
	}

	lastProductStock, _ := repository.GetLastProductStock(data.ID, []string{})
	response.Stock = lastProductStock.Current

	return
}

func CreateProduct(userID string, request dto.ProductRequest) (response models.SDAProduct, statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = errors.New("failed to parse user UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}
	if request.Code == "" || request.Code == "-" {
		request.Code = utils.GenerateRandomNumber(12)
	}

	parsedCategoryUUID, err := uuid.Parse(request.CategoryID)
	if err != nil {
		err = errors.New("failed to parse category UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	jsonImage, err := json.Marshal(request.Image)
	if err != nil {
		err = errors.New("failed to marshal image: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data := models.SDAProduct{
		Code:        request.Code,
		Name:        request.Name,
		Description: request.Description,
		Status:      request.Status,
		UserID:      parsedUserUUID,
		CategoryID:  parsedCategoryUUID,
		Price:       request.Price,
		Image:       string(jsonImage),
	}

	response, err = repository.CreateProduct(data)
	if err != nil {
		err = errors.New("failed to create data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusCreated
	return
}

func GetProductByID(id string, preloadFields []string) (response dto.ProductResponse, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = errors.New("failed to parse UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data, err := repository.GetProductByID(parsedUUID, preloadFields)
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	response, err = BuildProductResponse(data)
	if err != nil {
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func GetProducts(name, userID, categoryID string, param utils.PagingRequest, preloadFields []string) (response utils.PagingResponse, data []models.SDAProduct, statusCode int, err error) {
	baseFilter := "deleted_at IS NULL"
	var baseFilterValues []any
	if userID != "" {
		baseFilter += " AND user_id = ?"
		baseFilterValues = append(baseFilterValues, userID)
	}
	filter := baseFilter
	filterValues := baseFilterValues

	if name != "" {
		filter += " AND name = ?"
		filterValues = append(filterValues, name)
	}
	if categoryID != "" {
		filter += " AND category_id = ?"
		filterValues = append(filterValues, categoryID)
	}
	if param.Custom != "" {
		filter += " AND status = ?"
		filterValues = append(filterValues, param.Custom.(string))
	}
	if param.Search != "" {
		filter += " AND (name ILIKE ? OR description ILIKE ?)"
		filterValues = append(filterValues, fmt.Sprintf("%%%s%%", param.Search))
		filterValues = append(filterValues, fmt.Sprintf("%%%s%%", param.Search))
	}

	data, total, totalFiltered, err := repository.GetProducts(dto.FindParameter{
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
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	var results []dto.ProductResponse
	for _, item := range data {
		var result dto.ProductResponse
		result, err = BuildProductResponse(item)
		if err != nil {
			statusCode = http.StatusInternalServerError
			return
		}

		results = append(results, result)
	}

	response = utils.PopulateResPaging(&param, results, total, totalFiltered)

	statusCode = http.StatusOK
	return
}

func UpdateProduct(id string, request dto.ProductRequest) (response models.SDAProduct, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = errors.New("failed to parse UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data, err := repository.GetProductByID(parsedUUID, []string{})
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
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
	if request.Price > 0 {
		data.Price = request.Price
	}

	var parsedCategoryUUID uuid.UUID
	if request.CategoryID != "" {
		parsedCategoryUUID, err = uuid.Parse(request.CategoryID)
		if err != nil {
			err = errors.New("failed to parse category UUID: " + err.Error())
			statusCode = http.StatusInternalServerError
			return
		}
		data.CategoryID = parsedCategoryUUID
	}

	data.Status = request.Status

	var jsonImage []byte
	if len(request.Image) > 0 {
		jsonImage, err = json.Marshal(request.Image)
		if err != nil {
			err = errors.New("failed to marshal image: " + err.Error())
			statusCode = http.StatusInternalServerError
			return
		}

		data.Image = string(jsonImage)
	}

	response, err = repository.UpdateProduct(data)
	if err != nil {
		err = errors.New("failed to update data: " + err.Error())
		statusCode = http.StatusInternalServerError
	}

	statusCode = http.StatusOK
	return
}

func DeleteProduct(id string) (statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = errors.New("failed to parse UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data, err := repository.GetProductByID(parsedUUID, []string{})
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	err = repository.DeleteProduct(data)
	if err != nil {
		err = errors.New("failed to delete data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func UploadProductPicture(file *multipart.FileHeader, userID string) (responseURL string, statusCode int, err error) {
	extension := filepath.Ext(file.Filename)
	if extension != ".png" && extension != ".jpg" && extension != ".jpeg" && extension != ".webp" {
		err = errors.New("the file extension is wrong. allowed file extensions are images (.png, .jpg, .jpeg, .webp)")
		statusCode = http.StatusBadRequest
		return
	}

	var src multipart.File
	src, err = file.Open()
	if err != nil {
		err = errors.New("faield to open file: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}
	defer src.Close()

	responseURL, _, _, err = upload.UploadFile(src, userID, "")
	if err != nil {
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func ImportProduct(file *multipart.FileHeader, userID string) (responses []models.SDAProduct, statusCode int, err error) {
	rows, err := utils.ValidateImportFile(file, 6)
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
		var response models.SDAProduct

		var category models.SDAProductCategory
		existingCategories, _, _, _ := repository.GetProductCategories(dto.FindParameter{
			Filter:       "deleted_at IS NULL AND code = ?",
			FilterValues: []any{data[5]},
		}, []string{})

		if len(existingCategories) > 0 {
			category = existingCategories[0]
		} else {
			existingCategories, _, _, _ = repository.GetProductCategories(dto.FindParameter{
				Filter:       "deleted_at IS NULL AND name = ?",
				FilterValues: []any{data[4]},
			}, []string{})

			if len(existingCategories) > 0 {
				category = existingCategories[0]
			} else {
				input := dto.ProductCategoryRequest{
					Code:   data[5],
					Name:   data[4],
					Status: true,
				}
				if input.Code == "" || input.Code == "-" {
					input.Code = utils.GenerateRandomNumber(12)
				}

				category, statusCode, err = CreateProductCategory(userID, input)
				if err != nil {
					return
				}
			}
		}

		check, _, _, _ := repository.GetProducts(dto.FindParameter{
			Filter:       "deleted_at IS NULL AND code = ?",
			FilterValues: []any{data[2]},
		}, []string{})

		if len(check) > 0 {
			response = check[0]
		} else {
			check, _, _, _ = repository.GetProducts(dto.FindParameter{
				Filter:       "deleted_at IS NULL AND name = ?",
				FilterValues: []any{data[0]},
			}, []string{})

			if len(check) > 0 {
				response = check[0]
			}
		}

		if response.ID == uuid.Nil {
			input := dto.ProductRequest{
				Code:        data[2],
				Name:        data[0],
				Description: data[1],
				Status:      true,
				Image:       []string{},
				CategoryID:  category.ID.String(),
			}
			price, _ := strconv.ParseFloat(data[3], 64)
			input.Price = price
			if input.Code == "" || input.Code == "-" {
				input.Code = utils.GenerateRandomNumber(12)
			}

			response, statusCode, err = CreateProduct(userID, input)
			if err != nil {
				return
			}
		}

		responses = append(responses, response)
	}

	return
}

func ExportProduct(userID, fileExtentison string) (remoteFile bytes.Buffer, filename string, statusCode int, err error) {
	filename = fmt.Sprintf("products_%v.%v", userID, fileExtentison)
	directory := fmt.Sprintf("/assets/download/%v", filename)

	products, _, _, err := repository.GetProducts(dto.FindParameter{
		BaseFilter:       "deleted_at IS NULL AND user_id = ?",
		BaseFilterValues: []any{userID},
	}, []string{})
	if err != nil {
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	var file bytes.Buffer

	if fileExtentison == "xlsx" {
		f := excelize.NewFile()
		sheetName := "Products"
		f.SetSheetName("Sheet1", sheetName)

		f.SetCellValue(sheetName, "A1", "Name")
		f.SetCellValue(sheetName, "B1", "Description")
		f.SetCellValue(sheetName, "C1", "Code")
		f.SetCellValue(sheetName, "D1", "Price")
		f.SetCellValue(sheetName, "E1", "Category")
		f.SetCellValue(sheetName, "F1", "Category Code")

		for i, item := range products {
			row := i + 2
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "-")
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "-")
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "-")
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), "-")
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), "-")
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), "-")

			if item.Name != "" {
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), item.Name)
			}
			if item.Description != "" {
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), item.Description)
			}
			if item.Code != "" {
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), item.Code)
			}
			if item.Price > 0 {
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), fmt.Sprintf("%.2f", item.Price))
			}
			if item.CategoryID != uuid.Nil {
				category, _ := repository.GetProductCategoryByID(item.CategoryID, []string{})
				if category.Name != "" {
					f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), category.Name)
				}
				if category.Code != "" {
					f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), category.Code)
				}
			}
		}

		err = f.Write(&file)
		if err != nil {
			statusCode = http.StatusInternalServerError
			err = errors.New("failed to save excel file: " + err.Error())
		}
	}

	if fileExtentison == "csv" {
		writer := csv.NewWriter(&file)
		defer writer.Flush()

		header := []string{"Name", "Description", "Code"}
		err = writer.Write(header)
		if err != nil {
			statusCode = http.StatusInternalServerError
			err = errors.New("failed to write header into csv file: " + err.Error())
			return
		}

		for _, item := range products {
			name := "-"
			description := "-"
			code := "-"
			price := "-"
			categoryName := "-"
			categoryCode := "-"

			if item.Name != "" {
				name = item.Name
			}
			if item.Description != "" {
				description = item.Description
			}
			if item.Code != "" {
				code = item.Code
			}
			if item.Price > 0 {
				price = fmt.Sprintf("%.2f", item.Price)
			}
			if item.CategoryID != uuid.Nil {
				category, _ := repository.GetProductCategoryByID(item.CategoryID, []string{})
				if category.Name != "" {
					categoryName = category.Name
				}
				if category.Code != "" {
					categoryCode = category.Code
				}
			}

			row := []string{name, description, code, price, categoryName, categoryCode}
			err = writer.Write(row)
			if err != nil {
				statusCode = http.StatusInternalServerError
				err = errors.New("failed to write data into csv file: " + err.Error())
				return
			}
		}
	}

	statusCode, err = upload.WriteRemoteFile(file, directory)
	if err != nil {
		statusCode = http.StatusInternalServerError
		err = errors.New("failed to write remote data: " + err.Error())
		return
	}

	remoteFile, statusCode, err = upload.GetRemoteFile(directory)

	return
}
