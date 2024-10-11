package repository

import (
	"log"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/google/uuid"
)

func CreateSale(data models.SDASale) (models.SDASale, error) {
	err := config.DB.Preload("User").Preload("Customer").Preload("Details").Create(&data).Error
	if err != nil {
		log.Printf("Failed to insert data to database: %v", err)
	}

	return data, err
}

func GetSaleByID(id uuid.UUID, preloadFields []string) (response models.SDASale, err error) {
	db := utils.BuildPreload(config.DB, preloadFields)

	err = db.Where("id = ?", id).First(&response).Error
	if err != nil {
		log.Printf("Failed to get data from database: %v", err)
	}

	return
}

func GetSales(param dto.FindParameter, preloadFields []string) (responses []models.SDASale, total int64, totalFiltered int64, err error) {
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

	db := utils.BuildPreload(config.DB, preloadFields)

	if param.Limit == 0 {
		err = db.Offset(param.Offset).Order(param.Order).Where(param.Filter).Find(&responses).Error
	} else {
		err = db.Limit(param.Limit).Offset(param.Offset).Order(param.Order).Where(param.Filter).Find(&responses).Error
	}
	if err != nil {
		log.Printf("Failed to get data list from database: %v", err)
	}

	return
}

func UpdateSale(data models.SDASale) (models.SDASale, error) {
	err := config.DB.Preload("User").Preload("Customer").Preload("Details").Save(&data).Error
	if err != nil {
		log.Printf("Failed to update data in database: %v", err)
	}

	return data, err
}

func DeleteSale(data models.SDASale) error {
	err := config.DB.Delete(&data).Error
	if err != nil {
		log.Printf("Failed to delete data in from database: %v", err)
	}

	return err
}
