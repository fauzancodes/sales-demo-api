package service

import (
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

	data := models.SDACustomer{
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
