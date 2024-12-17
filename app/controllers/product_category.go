package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/pkg/upload"
	"github.com/fauzancodes/sales-demo-api/app/pkg/utils"
	"github.com/fauzancodes/sales-demo-api/app/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func CreateProductCategory(c echo.Context) error {
	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	var request dto.ProductCategoryRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusUnprocessableEntity,
			dto.Response{
				Status:  http.StatusUnprocessableEntity,
				Message: "Invalid request body",
				Error:   err.Error(),
			},
		)
	}

	if err := request.Validate(); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Status:  http.StatusBadRequest,
				Message: "Invalid request value",
				Error:   err.Error(),
			},
		)
	}

	result, statusCode, err := service.CreateProductCategory(userID, request)
	if err != nil {
		return c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to create",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to create",
			Data:    result,
		},
	)
}

func GetProductCategories(c echo.Context) error {
	withUser, _ := strconv.ParseBool(c.QueryParam("with_user"))
	withProducts, _ := strconv.ParseBool(c.QueryParam("with_products"))

	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	name := c.QueryParam("name")

	var preloadFields []string
	if withUser {
		preloadFields = append(preloadFields, "User")
	}
	if withProducts {
		preloadFields = append(preloadFields, "Products")
	}

	param := utils.PopulatePaging(c, "status")
	data, _, statusCode, err := service.GetProductCategories(name, userID, param, preloadFields)
	if err != nil {
		return c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to get data",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(statusCode, data)
}

func GetProductCategoryByID(c echo.Context) error {
	withUser, _ := strconv.ParseBool(c.QueryParam("with_user"))
	withProducts, _ := strconv.ParseBool(c.QueryParam("with_products"))

	id := c.Param("id")

	var preloadFields []string
	if withUser {
		preloadFields = append(preloadFields, "User")
	}
	if withProducts {
		preloadFields = append(preloadFields, "Products")
	}

	data, statusCode, err := service.GetProductCategoryByID(id, preloadFields)
	if err != nil {
		return c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to get data",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to get data",
			Data:    data,
		},
	)
}

func UpdateProductCategory(c echo.Context) error {
	id := c.Param("id")

	var request dto.ProductCategoryRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(
			http.StatusUnprocessableEntity,
			dto.Response{
				Status:  http.StatusUnprocessableEntity,
				Message: "Invalid request body",
				Error:   err.Error(),
			},
		)
	}

	data, statusCode, err := service.UpdateProductCategory(id, request)
	if err != nil {
		return c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to update data",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to update data",
			Data:    data,
		},
	)
}

func DeleteProductCategory(c echo.Context) error {
	id := c.Param("id")

	statusCode, err := service.DeleteProductCategory(id)
	if err != nil {
		return c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to delete data",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to delete data",
		},
	)
}

func GetProductCategoryImportTemplate(c echo.Context) error {
	buf, statusCode, err := upload.GetRemoteFile("/assets/template/product_category.xlsx")
	if err != nil {
		return c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to get file",
				Error:   err.Error(),
			},
		)
	}

	c.Response().Header().Set("Content-Disposition", `attachment; filename="product_category.xlsx"`)
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	return c.Blob(http.StatusOK, "application/octet-stream", buf.Bytes())
}

func ImportProductCategory(c echo.Context) error {
	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	file, err := c.FormFile("file")
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

	response, statusCode, err := service.ImportProductCategory(file, userID)
	if err != nil {
		return c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to import data",
				Error:   err.Error(),
			},
		)
	}

	return c.JSON(
		statusCode,
		dto.Response{
			Status:  statusCode,
			Message: "Success to import data",
			Data:    response,
		},
	)
}

func ExportProductCategory(c echo.Context) error {
	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	fileExtension := c.QueryParam("file_extension")
	if fileExtension != "xlsx" && fileExtension != "csv" {
		return c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Status:  http.StatusBadRequest,
				Message: "the file format only accepts .xlsx and .csv",
			},
		)
	}

	remoteFile, filename, statusCode, err := service.ExportProductCategory(userID, fileExtension)
	if err != nil {
		return c.JSON(
			statusCode,
			dto.Response{
				Status:  statusCode,
				Message: "Failed to export product category",
				Error:   err.Error(),
			},
		)
	}

	c.Response().Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	c.Response().Header().Set("Content-Type", "application/octet-stream")

	return c.Blob(http.StatusOK, "application/octet-stream", remoteFile.Bytes())
}
