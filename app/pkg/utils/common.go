package utils

import (
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"
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

func GetBaseUrl(c echo.Context) (response string) {
	response = fmt.Sprintf("%v://%v", c.Scheme(), c.Request().Host)

	return
}

func ValidateImportFile(file *multipart.FileHeader, numberOfColumns int) (rows [][]string, err error) {
	src, err := file.Open()
	if err != nil {
		return
	}
	defer src.Close()

	extension := filepath.Ext(file.Filename)
	if extension == ".xls" || extension == ".xlsx" {
		var f *excelize.File
		f, err = excelize.OpenReader(src)
		if err != nil {
			return
		}

		sheets := f.GetSheetList()
		if len(sheets) == 0 {
			err = errors.New("there is no sheet in the file")
			return
		}

		rows, err = f.GetRows(sheets[0])
		if err != nil {
			return
		}
		if len(rows[0]) != numberOfColumns {
			err = errors.New("The number of columns must match the template. Expected: " + strconv.Itoa(numberOfColumns) + " columns. Current: " + strconv.Itoa(len(rows[0])) + " columns")
			return
		}
	} else if extension == ".csv" {
		reader := csv.NewReader(src)
		rows, err = reader.ReadAll()
		if err != nil {
			return
		}
	} else {
		err = errors.New("the file format only accepts .xls, .xlsx, .csv")
		return
	}

	return
}

func RoundFloat(number float64) (result int) {
	return int(math.Round(number))
}
