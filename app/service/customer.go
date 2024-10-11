package service

import (
	"errors"
	"mime/multipart"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/repository"
	"github.com/google/uuid"
)

func CreateCustomer(userID string, request dto.CustomerRequest) (response models.SDACustomer, err error) {
	parsedUserUUID, err := uuid.Parse(userID)
	if err != nil {
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

	return
}

func GetCustomerByID(id string, preloadFields []string) (data models.SDACustomer, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return
	}
	data, err = repository.GetCustomerByID(parsedUUID, preloadFields)

	return
}

func GetCustomers(email, phone, userID string, param utils.PagingRequest, preloadFields []string) (response utils.PagingResponse, data []models.SDACustomer, err error) {
	baseFilter := "deleted_at IS NULL"
	if userID != "" {
		baseFilter += " AND user_id = '" + userID + "'"
	}
	filter := baseFilter

	if email != "" {
		filter += " AND email = '" + email + "'"
	}
	if phone != "" {
		filter += " AND phone = '" + phone + "'"
	}
	if param.Custom.(string) != "" {
		filter += " AND status = " + param.Custom.(string)
	}
	if param.Search != "" {
		filter += " AND (first_name ILIKE '%" + param.Search + "%' OR last_name ILIKE '%" + param.Search + "%')"
	}

	data, total, totalFiltered, err := repository.GetCustomers(dto.FindParameter{
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

func UpdateCustomer(id string, request dto.CustomerRequest) (response models.SDACustomer, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return
	}

	data, err := repository.GetCustomerByID(parsedUUID, []string{})
	if err != nil {
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

	return
}

func DeleteCustomer(id string) (err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return
	}

	data, err := repository.GetCustomerByID(parsedUUID, []string{})
	if err != nil {
		return
	}

	err = repository.DeleteCustomer(data)

	return
}

func ImportCustomer(file *multipart.FileHeader, userID string) (responses []models.SDACustomer, err error) {
	rows, err := utils.ValidateImportFile(file, 5)
	if err != nil {
		return
	}

	rows = rows[1:]
	if len(rows) > 0 {
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

				response, err = CreateCustomer(userID, input)
				if err != nil {
					return
				}
			}

			responses = append(responses, response)
		}
	} else {
		err = errors.New("there is no data in the file")
		return
	}

	return
}
