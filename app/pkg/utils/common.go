package utils

import (
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func BuildPreload(db *gorm.DB, fields []string) *gorm.DB {
	if len(fields) > 0 {
		for _, field := range fields {
			db = db.Preload(field)
		}
	}

	return db
}

func GetBuildPreloadFields(c echo.Context) (fields []string) {
	raw := c.QueryParam("preload_fields")

	fields = strings.Split(raw, ",")

	return
}
