package repository

import (
	"log"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/google/uuid"
)

func CreateUser(data models.SDAUser) (models.SDAUser, error) {
	err := config.DB.Create(&data).Error
	if err != nil {
		log.Printf("Failed to insert data to database: %v", err)
	}

	return data, err
}

func GetUserByID(id uuid.UUID) (response models.SDAUser, err error) {
	err = config.DB.First(&response, id).Error
	if err != nil {
		log.Printf("Failed to get data from database: %v", err)
	}

	return
}

func GetUsers(param dto.FindParameter) (responses []models.SDAUser, total int64, totalFiltered int64, err error) {
	err = config.DB.Model(responses).Where(param.BaseFilter).Count(&total).Error
	if err != nil {
		log.Printf("Failed to count data from database: %v", err)
		return
	}

	err = config.DB.Model(responses).Where(param.Filter).Count(&totalFiltered).Error
	if err != nil {
		log.Printf("Failed to count filtered data from database: %v", err)
		return
	}

	if param.Limit == 0 {
		err = config.DB.Offset(param.Offset).Order(param.Order).Where(param.Filter).Find(&responses).Error
	} else {
		err = config.DB.Limit(param.Limit).Offset(param.Offset).Order(param.Order).Where(param.Filter).Find(&responses).Error
	}
	if err != nil {
		log.Printf("Failed to get data list from database: %v", err)
	}

	return
}

func UpdateUser(data models.SDAUser) (models.SDAUser, error) {
	err := config.DB.Save(&data).Error
	if err != nil {
		log.Printf("Failed to update data in database: %v", err)
	}

	return data, err
}

func DeleteUser(data models.SDAUser) (models.SDAUser, error) {
	err := config.DB.Delete(&data).Error
	if err != nil {
		log.Printf("Failed to delete data in from database: %v", err)
	}

	return data, err
}
