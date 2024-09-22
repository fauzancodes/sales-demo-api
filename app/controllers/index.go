package controllers

import (
	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
	return c.JSON(200, "Welcome to Sales Demo API")
}
