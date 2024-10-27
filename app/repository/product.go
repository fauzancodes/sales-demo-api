package repository

import (
	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/google/uuid"
)

func CreateProduct(data models.SDAProduct) (models.SDAProduct, error) {
	err := config.DB.Create(&data).Error

	return data, err
}

func GetProductByID(id uuid.UUID, preloadFields []string) (response models.SDAProduct, err error) {
	db := utils.BuildPreload(config.DB, preloadFields)

	err = db.Where("id = ?", id).First(&response).Error

	return
}

func GetProducts(param dto.FindParameter, preloadFields []string) (responses []models.SDAProduct, total int64, totalFiltered int64, err error) {
	err = config.DB.Model(responses).Where(param.BaseFilter).Count(&total).Error
	if err != nil {
		return
	}

	err = config.DB.Model(responses).Where(param.Filter).Count(&totalFiltered).Error
	if err != nil {
		return
	}

	db := utils.BuildPreload(config.DB, preloadFields)

	if param.Limit == 0 {
		err = db.Offset(param.Offset).Order(param.Order).Where(param.Filter).Find(&responses).Error
	} else {
		err = db.Limit(param.Limit).Offset(param.Offset).Order(param.Order).Where(param.Filter).Find(&responses).Error
	}

	return
}

func UpdateProduct(data models.SDAProduct) (models.SDAProduct, error) {
	err := config.DB.Save(&data).Error

	return data, err
}

func DeleteProduct(data models.SDAProduct) error {
	err := config.DB.Delete(&data).Error

	return err
}
