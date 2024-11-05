package repository

import (
	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/google/uuid"
)

func CreateProductStock(data models.SDAProductStock) (models.SDAProductStock, error) {
	err := config.DB.Create(&data).Error

	return data, err
}

func GetProductStocks(param dto.FindParameter, preloadFields []string) (responses []models.SDAProductStock, total int64, totalFiltered int64, err error) {
	err = config.DB.Model(responses).Where(param.BaseFilter, param.BaseFilterValues...).Count(&total).Error
	if err != nil {
		return
	}

	err = config.DB.Model(responses).Where(param.Filter, param.FilterValues...).Count(&totalFiltered).Error
	if err != nil {
		return
	}

	db := utils.BuildPreload(config.DB, preloadFields)

	if param.Limit == 0 {
		err = db.Offset(param.Offset).Order(param.Order).Where(param.Filter, param.FilterValues...).Find(&responses).Error
	} else {
		err = db.Limit(param.Limit).Offset(param.Offset).Order(param.Order).Where(param.Filter, param.FilterValues...).Find(&responses).Error
	}

	return
}

func GetLastProductStock(id uuid.UUID, preloadFields []string) (response models.SDAProductStock, err error) {
	db := utils.BuildPreload(config.DB, preloadFields)

	err = db.Where("product_id = ?", id).Order("created_at DESC").First(&response).Error

	return
}
