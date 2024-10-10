package controllers

import (
	"log"
	"net/http"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func CreateProduct(c echo.Context) error {
	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	var request dto.ProductRequest
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

	result, err := service.CreateProduct(userID, request)
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

func GetProducts(c echo.Context) error {
	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	name := c.QueryParam("name")
	categoryID := c.QueryParam("category_id")
	preloadFields := utils.GetBuildPreloadFields(c)

	param := utils.PopulatePaging(c, "status")
	data, _, err := service.GetProducts(name, userID, categoryID, param, preloadFields)
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

func GetProductByID(c echo.Context) error {
	id := c.Param("id")
	preloadFields := utils.GetBuildPreloadFields(c)

	data, err := service.GetProductByID(id, preloadFields)
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

func UpdateProduct(c echo.Context) error {
	id := c.Param("id")

	var request dto.ProductRequest
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

	data, err := service.UpdateProduct(id, request)
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

func DeleteProduct(c echo.Context) error {
	id := c.Param("id")

	err := service.DeleteProduct(id)
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

func UploadProductPicture(c echo.Context) error {
	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Status:  500,
				Message: "Failed to get file from form",
				Error:   err.Error(),
			},
		)
	}

	responseURL, err := service.UploadProductPicture(file, userID)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			dto.Response{
				Status:  500,
				Message: "Failed to upload product picture",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		dto.Response{
			Status:  200,
			Message: "Success to upload",
			Data:    responseURL,
		},
	)
}
