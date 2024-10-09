package controllers

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/pkg/upload"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func UploadFile(c echo.Context) error {
	userID := c.Get("currentUser").(jwt.MapClaims)["id"].(string)
	log.Printf("Current user ID: %v", userID)

	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			dto.Response{
				Status:  500,
				Message: "Failed to initiate form",
				Error:   err.Error(),
			},
		)
	}

	file := form.File["image"]
	if len(file) == 0 {
		return c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Status:  500,
				Message: "No files selected",
				Error:   "bad request",
			},
		)
	}

	var responseURL string
	extension := filepath.Ext(file[0].Filename)
	if extension == ".png" || extension == ".jpg" || extension == ".jpeg" || extension == ".webp" {
		src, err := file[0].Open()
		if err != nil {
			return c.JSON(
				http.StatusBadRequest,
				dto.Response{
					Status:  500,
					Message: "Failed to open file",
					Error:   err.Error(),
				},
			)
		}
		defer src.Close()

		secureUrl, _, _, err := upload.UploadFile(src, userID, "")
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				dto.Response{
					Status:  500,
					Message: "Failed to upload file",
					Error:   err.Error(),
				},
			)
		}

		responseURL = secureUrl
	} else {
		return c.JSON(
			http.StatusBadRequest,
			dto.Response{
				Status:  500,
				Message: "The file extension is wrong. Allowed file extensions are images (.png, .jpg, .jpeg, .webp)",
				Error:   "bad request",
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
