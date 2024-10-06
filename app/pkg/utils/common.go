package utils

import (
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/exp/rand"
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

	if raw != "" {
		fields = strings.Split(raw, ",")
	}

	return
}

func GenerateRandomNumber(length int) string {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		location = time.Local
		err = nil
	}
	rand.Seed(uint64(time.Now().In(location).UnixNano()))
	charset := "0123456789"
	randomBytes := make([]byte, length)
	for i := range randomBytes {
		randomBytes[i] = charset[rand.Intn(len(charset))]
	}
	randomString := string(randomBytes)
	return randomString
}
