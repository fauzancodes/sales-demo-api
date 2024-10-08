package controllers

import (
	"log"
	"net/http"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func CreateSale(c echo.Context) error {
	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	var request dto.SaleRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusUnprocessableEntity,
			dto.Response{
				Status:  422,
				Message: "Invalid request body",
				Error:   err.Error(),
			},
		)
	}

	if err := request.Validate(); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Status:  400,
				Message: "Invalid request value",
				Error:   err.Error(),
			},
		)
	}

	result, err := service.CreateSale(userID, request)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			dto.Response{
				Status:  500,
				Message: "Failed to create",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  200,
			Message: "Success to create",
			Data:    result,
		},
	)
}

func GetSales(c echo.Context) error {
	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	invoiceID := c.QueryParam("invoice_id")
	transactionDateMarginTop := c.QueryParam("transaction_date_margin_top")
	transactionDateMarginBottom := c.QueryParam("transaction_date_margin_bottom")
	productID := c.QueryParam("product_id")
	customerID := c.QueryParam("customer_id")
	preloadFields := utils.GetBuildPreloadFields(c)

	param := utils.PopulatePaging(c, "status")
	data, _, err := service.GetSales(invoiceID, userID, customerID, transactionDateMarginTop, transactionDateMarginBottom, productID, param, preloadFields)
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

func GetSaleByID(c echo.Context) error {
	id := c.Param("id")
	preloadFields := utils.GetBuildPreloadFields(c)

	data, err := service.GetSaleByID(id, preloadFields)
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

	return c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  200,
			Message: "Success to get data",
			Data:    data,
		},
	)
}

func UpdateSale(c echo.Context) error {
	id := c.Param("id")

	var request dto.SaleRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusUnprocessableEntity,
			dto.Response{
				Status:  422,
				Message: "Invalid request body",
				Error:   err.Error(),
			},
		)
	}

	data, err := service.UpdateSale(id, request)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			dto.Response{
				Status:  500,
				Message: "Failed to update data",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  200,
			Message: "Success to update data",
			Data:    data,
		},
	)
}

func DeleteSale(c echo.Context) error {
	id := c.Param("id")

	err := service.DeleteSale(id)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			dto.Response{
				Status:  500,
				Message: "Failed to delete data",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  200,
			Message: "Success to delete data",
		},
	)
}

func SendSaleInvoiceByID(c echo.Context) error {
	id := c.Param("id")
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Status:  500,
				Message: "Failed to parse uuid",
				Error:   err.Error(),
			},
		)
	}

	go service.SendSaleInvoice(parsedUUID)

	return c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  200,
			Message: "Success to send sale invoice to email",
		},
	)
}
