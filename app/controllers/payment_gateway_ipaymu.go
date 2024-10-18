package controllers

import (
	"net/http"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/service"
	"github.com/labstack/echo/v4"
)

func GetIPaymuPaymentMethods(c echo.Context) error {
	code := c.QueryParam("code")

	param := utils.PopulatePaging(c, "")
	data, _, err := service.GetIPaymuPaymentMethods(code, param)
	if err != nil {
		return c.JSON(
			http.StatusNotFound,
			dto.Response{
				Status:  404,
				Message: "Failed to get data",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(http.StatusOK, data)
}
