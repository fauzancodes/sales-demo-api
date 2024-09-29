package repository

import (
	"log"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/google/uuid"
)

func CreateProductCategory(data models.SDAProductCategory) (models.SDAProductCategory, error) {
	err := config.DB.Create(&data).Error
	if err != nil {
		log.Printf("Failed to insert data to database: %v", err)
	}

	return data, err
}

func GetProductCategoryByID(id uuid.UUID) (response models.SDAProductCategory, err error) {
	err = config.DB.Where("id = ?", id).First(&response).Error
	if err != nil {
		log.Printf("Failed to get data from database: %v", err)
	}

	return
}

func GetProductCategories(param dto.FindParameter) (responses []models.SDAProductCategory, total int64, totalFiltered int64, err error) {
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

func UpdateProductCategory(data models.SDAProductCategory) (models.SDAProductCategory, error) {
	err := config.DB.Save(&data).Error
	if err != nil {
		log.Printf("Failed to update data in database: %v", err)
	}

	return data, err
}

func DeleteProductCategory(data models.SDAProductCategory) error {
	err := config.DB.Delete(&data).Error
	if err != nil {
		log.Printf("Failed to delete data in from database: %v", err)
	}

	return err
}