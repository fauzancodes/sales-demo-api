package service

import (
	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/repository"
)

func GetIPaymuPaymentMethods(code string, param utils.PagingRequest) (response utils.PagingResponse, data []models.SDAIPaymuPaymentMethod, err error) {
	baseFilter := "deleted_at IS NULL"
	filter := baseFilter

	if code != "" {
		filter += " AND code = '" + code + "'"
	}
	if param.Search != "" {
		filter += " AND (name ILIKE '%" + param.Search + "%' OR description ILIKE '%" + param.Search + "%')"
	}

	data, total, totalFiltered, err := repository.GetIPaymuPaymentMethods(dto.FindParameter{
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
