package repository

import (
	"log"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/google/uuid"
)

func CreateCustomer(data models.SDACustomer) (models.SDACustomer, error) {
	err := config.DB.Create(&data).Error
	if err != nil {
		log.Printf("Failed to insert data to database: %v", err)
	}

	return data, err
}

func GetCustomerByID(id uuid.UUID) (response models.SDACustomer, err error) {
	err = config.DB.Where("id = ?", id).First(&response).Error
	if err != nil {
		log.Printf("Failed to get data from database: %v", err)
	}

	return
}

func GetCustomers(param dto.FindParameter) (responses []models.SDACustomer, total int64, totalFiltered int64, err error) {
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

func UpdateCustomer(data models.SDACustomer) (models.SDACustomer, error) {
	err := config.DB.Save(&data).Error
	if err != nil {
		log.Printf("Failed to update data in database: %v", err)
	}

	return data, err
}

func DeleteCustomer(data models.SDACustomer) error {
	err := config.DB.Delete(&data).Error
	if err != nil {
		log.Printf("Failed to delete data in from database: %v", err)
	}

	return err
}