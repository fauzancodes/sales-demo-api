package service

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
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

func CreateCustomer(userID string, request dto.CustomerRequest) (response models.SDACustomer, statusCode int, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
		err = errors.New("failed to parse user UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}
	if request.Code == "" || request.Code == "-" {
		request.Code = utils.GenerateRandomNumber(12)
	}

	data := models.SDACustomer{
		Code:      request.Code,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Phone:     request.Phone,
		Status:    request.Status,
		UserID:    parsedUserUUID,
	}

	response, err = repository.CreateCustomer(data)
	if err != nil {
		err = errors.New("failed to create data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusCreated
	return
}

func GetCustomerByID(id string, preloadFields []string) (data models.SDACustomer, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = errors.New("failed to parse UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data, err = repository.GetCustomerByID(parsedUUID, preloadFields)
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func GetCustomers(email, phone, userID string, param utils.PagingRequest, preloadFields []string) (response utils.PagingResponse, data []models.SDACustomer, statusCode int, err error) {
	baseFilter := "deleted_at IS NULL"
	var baseFilterValues []any
	if userID != "" {
		baseFilter += " AND user_id = ?"
		baseFilterValues = append(baseFilterValues, userID)
	}
	filter := baseFilter
	filterValues := baseFilterValues

	if email != "" {
		filter += " AND email = ?"
		filterValues = append(filterValues, email)
	}
	if phone != "" {
		filter += " AND phone = ?"
		filterValues = append(filterValues, phone)
	}
	if param.Custom != "" {
		filter += " AND status = ?"
		filterValues = append(filterValues, param.Custom.(string))
	}
	if param.Search != "" {
		filter += " AND (first_name ILIKE ? OR last_name ILIKE ?)"
		filterValues = append(filterValues, fmt.Sprintf("%%%s%%", param.Search))
		filterValues = append(filterValues, fmt.Sprintf("%%%s%%", param.Search))
	}

	data, total, totalFiltered, err := repository.GetCustomers(dto.FindParameter{
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

	response = utils.PopulateResPaging(&param, data, total, totalFiltered)

	statusCode = http.StatusOK
	return
}

func UpdateCustomer(id string, request dto.CustomerRequest) (response models.SDACustomer, statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = errors.New("failed to parse UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data, err := repository.GetCustomerByID(parsedUUID, []string{})
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
	if request.FirstName != "" {
		data.FirstName = request.FirstName
	}
	if request.LastName != "" {
		data.LastName = request.LastName
	}
	if request.Email != "" {
		data.Email = request.Email
	}
	if request.Phone != "" {
		data.Phone = request.Phone
	}
	data.Status = request.Status

	response, err = repository.UpdateCustomer(data)
	if err != nil {
		err = errors.New("failed to update data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func DeleteCustomer(id string) (statusCode int, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		err = errors.New("failed to parse UUID: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	data, err := repository.GetCustomerByID(parsedUUID, []string{})
	if err != nil {
		err = errors.New("failed to get data: " + err.Error())
		if strings.Contains(err.Error(), gorm.ErrRecordNotFound.Error()) {
			statusCode = http.StatusNotFound
			return
		}

		statusCode = http.StatusInternalServerError
		return
	}

	err = repository.DeleteCustomer(data)
	if err != nil {
		err = errors.New("failed to delete data: " + err.Error())
		statusCode = http.StatusInternalServerError
		return
	}

	statusCode = http.StatusOK
	return
}

func ImportCustomer(file *multipart.FileHeader, userID string) (responses []models.SDACustomer, statusCode int, err error) {
	rows, err := utils.ValidateImportFile(file, 5)
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
		var response models.SDACustomer

		check, _, _, _ := repository.GetCustomers(dto.FindParameter{
			Filter: "deleted_at IS NULL AND code = '" + data[4] + "'",
		}, []string{})

		if len(check) > 0 {
			response = check[0]
		} else {
			check, _, _, _ = repository.GetCustomers(dto.FindParameter{
				Filter: "deleted_at IS NULL AND email = '" + data[2] + "'",
			}, []string{})

			if len(check) > 0 {
				response = check[0]
			}
		}

		if response.ID == uuid.Nil {
			input := dto.CustomerRequest{
				Code:      data[4],
				FirstName: data[0],
				LastName:  data[1],
				Email:     data[2],
				Phone:     data[3],
				Status:    true,
			}
			if input.Code == "" || input.Code == "-" {
				input.Code = utils.GenerateRandomNumber(12)
			}

			response, statusCode, err = CreateCustomer(userID, input)
			if err != nil {
				return
			}
		}

		responses = append(responses, response)
	}

	return
}

func ExportCustomer(userID, fileExtentison string) (remoteFile bytes.Buffer, filename string, statusCode int, err error) {
	filename = fmt.Sprintf("customers_%v.%v", userID, fileExtentison)
	directory := fmt.Sprintf("/assets/download/%v", filename)

	customers, _, _, err := repository.GetCustomers(dto.FindParameter{
		BaseFilter: "deleted_at IS NULL AND user_id = '" + userID + "'",
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
		sheetName := "Customers"
		f.SetSheetName("Sheet1", sheetName)

		f.SetCellValue(sheetName, "A1", "First Name")
		f.SetCellValue(sheetName, "B1", "Last Name")
		f.SetCellValue(sheetName, "C1", "Email")
		f.SetCellValue(sheetName, "D1", "Phone")
		f.SetCellValue(sheetName, "E1", "Code")

		for i, item := range customers {
			row := i + 2
			f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "-")
			f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "-")
			f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), "-")
			f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), "-")
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), "-")

			if item.FirstName != "" {
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), item.FirstName)
			}
			if item.LastName != "" {
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), item.LastName)
			}
			if item.Email != "" {
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), item.Email)
			}
			if item.Phone != "" {
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), item.Phone)
			}
			if item.Code != "" {
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), item.Code)
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

		header := []string{"First Name", "Last Name", "Email", "Phone", "Code"}
		err = writer.Write(header)
		if err != nil {
			statusCode = http.StatusInternalServerError
			err = errors.New("failed to write header into csv file: " + err.Error())
			return
		}

		for _, item := range customers {
			firstName := "-"
			lastName := "-"
			email := "-"
			phone := "-"
			code := "-"

			if item.FirstName != "" {
				firstName = item.FirstName
			}
			if item.LastName != "" {
				lastName = item.LastName
			}
			if item.Email != "" {
				email = item.Email
			}
			if item.Phone != "" {
				phone = item.Phone
			}
			if item.Code != "" {
				code = item.Code
			}

			row := []string{firstName, lastName, email, phone, code}
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
