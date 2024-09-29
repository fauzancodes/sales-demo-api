package controllers

import (
	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
	return c.File("assets/html/index.html")
}
