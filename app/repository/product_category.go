package repository

import (
	"errors"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/google/uuid"
)

func CreateProductCategory(data models.SDAProductCategory) (models.SDAProductCategory, error) {
	err := config.DB.Create(&data).Error
	if err != nil {
		err = errors.New("failed to insert data to database: " + err.Error())
	}

	return data, err
}

func GetProductCategoryByID(id uuid.UUID, preloadFields []string) (response models.SDAProductCategory, err error) {
	db := utils.BuildPreload(config.DB, preloadFields)

	err = db.Where("id = ?", id).First(&response).Error
	if err != nil {
		err = errors.New("failed to get data from database: " + err.Error())
	}

	return
}

func GetProductCategories(param dto.FindParameter, preloadFields []string) (responses []models.SDAProductCategory, total int64, totalFiltered int64, err error) {
	err = config.DB.Model(responses).Where(param.BaseFilter).Count(&total).Error
	if err != nil {
		err = errors.New("failed to count data from database: " + err.Error())
		return
	}

	err = config.DB.Model(responses).Where(param.Filter).Count(&totalFiltered).Error
	if err != nil {
		err = errors.New("failed to count filtered data from database: " + err.Error())
		return
	}

	db := utils.BuildPreload(config.DB, preloadFields)

	if param.Limit == 0 {
		err = db.Offset(param.Offset).Order(param.Order).Where(param.Filter).Find(&responses).Error
	} else {
		err = db.Limit(param.Limit).Offset(param.Offset).Order(param.Order).Where(param.Filter).Find(&responses).Error
	}
	if err != nil {
		err = errors.New("failed to get data list from database: " + err.Error())
	}

	return
}

func UpdateProductCategory(data models.SDAProductCategory) (models.SDAProductCategory, error) {
	err := config.DB.Save(&data).Error
	if err != nil {
		err = errors.New("failed to update data in database: " + err.Error())
	}

	return data, err
}

func DeleteProductCategory(data models.SDAProductCategory) error {
	err := config.DB.Delete(&data).Error
	if err != nil {
		err = errors.New("failed to delete data in from database: " + err.Error())
	}

	return err
}
