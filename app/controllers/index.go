package controllers

import (
	"net/http"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/pkg/upload"
	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
	buf, statusCode, err := upload.GetRemoteFile("/assets/html/index.html")
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

	return c.Blob(http.StatusOK, "text/html", buf.Bytes())
}

func DownloadPostmanCollection(c echo.Context) error {
	buf, statusCode, err := upload.GetRemoteFile("/docs/Sales Demo API.postman_collection.json")
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

	c.Response().Header().Set("Content-Disposition", `attachment; filename="Sales Demo API.postman_collection.json"`)
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	return c.Blob(http.StatusOK, "application/octet-stream", buf.Bytes())
}

func DownloadPostmanEnvironment(c echo.Context) error {
	buf, statusCode, err := upload.GetRemoteFile("/docs/Sales Demo API.postman_environment.json")
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

	c.Response().Header().Set("Content-Disposition", `attachment; filename="Sales Demo API.postman_environment.json"`)
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	return c.Blob(http.StatusOK, "application/octet-stream", buf.Bytes())
}
