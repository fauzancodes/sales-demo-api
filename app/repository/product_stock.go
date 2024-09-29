package repository

import (
	"log"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/google/uuid"
)

func CreateProductStock(data models.SDAProductStock) (models.SDAProductStock, error) {
	err := config.DB.Create(&data).Error
	if err != nil {
		log.Printf("Failed to insert data to database: %v", err)
	}

	return data, err
}

func GetProductStocks(param dto.FindParameter) (responses []models.SDAProductStock, total int64, totalFiltered int64, err error) {
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

func GetLastProductStock(id uuid.UUID) (response models.SDAProductStock, err error) {
	err = config.DB.Where("product_id = ?", id).Order("created_at DESC").First(&response).Error
	if err != nil {
		log.Printf("Failed to get data from database: %v", err)
	}

	return
}
