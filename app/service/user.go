package service

import (
	"log"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/repository"
	"github.com/fauzancodes/sales-demo-api/pkg/bcrypt"
	"github.com/fauzancodes/sales-demo-api/pkg/utils"
	"github.com/google/uuid"
)

func CreateUser(request dto.UserRequest) (response models.SDAUser, err error) {
	data := models.SDAUser{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  bcrypt.HashPassword(request.Password),
	}

	response, err = repository.CreateUser(data)

	return
}

func GetUserByID(id string) (response models.SDAUser, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Failed to parse UUID: %v", err)
		return
	}

	response, err = repository.GetUserByID(parsedUUID)

	return
}

func GetUsers(firstName, lastName, email string, param utils.PagingRequest) (response utils.PagingResponse, data []models.SDAUser, err error) {
	baseFilter := "deleted_at IS NULL"
	filter := baseFilter

	if firstName != "" {
		filter += " AND first_name = '" + firstName + "'"
	}
	if lastName != "" {
		filter += " AND last_name = '" + lastName + "'"
	}
	if email != "" {
		filter += " AND email = '" + email + "'"
	}
	if param.Search != "" {
		filter += " AND (first_name ILIKE '%" + param.Search + "%' OR last_name ILIKE '%" + param.Search + "%' OR email ILIKE '%" + param.Search + "%')"
	}

	data, total, totalFiltered, err := repository.GetUsers(dto.FindParameter{
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

func UpdateUser(id string) (response models.SDAUser, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Failed to parse UUID: %v", err)
		return
	}

	data, err := repository.GetUserByID(parsedUUID)
	if err != nil {
		return
	}

	response, err = repository.UpdateUser(data)

	return
}

func DeleteUser(id string) (response models.SDAUser, err error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		log.Printf("Failed to parse UUID: %v", err)
		return
	}

	data, err := repository.GetUserByID(parsedUUID)
	if err != nil {
		return
	}

	response, err = repository.DeleteUser(data)

	return
}