package repository

import (
	"errors"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
)

func GetXenditPaymentMethods(param dto.FindParameter) (responses []models.SDAXenditPaymentMethod, total int64, totalFiltered int64, err error) {
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

	db := config.DB

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

func CreateXenditSalePayment(data models.SDAXenditSalePayment) (models.SDAXenditSalePayment, error) {
	err := config.DB.Create(&data).Error
	if err != nil {
		err = errors.New("failed to insert data to database: " + err.Error())
		return data, err
	}

	err = config.DB.Preload("PaymentMethod").First(&data, data.ID).Error
	if err != nil {
		err = errors.New("failed to load PaymentMethod: " + err.Error())
		return data, err
	}

	return data, err
}

func GetXenditSalePayments(param dto.FindParameter) (responses []models.SDAXenditSalePayment, total int64, totalFiltered int64, err error) {
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

	db := utils.BuildPreload(config.DB, []string{"PaymentMethod"})

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
