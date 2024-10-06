package config

import (
	"os"
	"strconv"
)

type Config struct {
	IndexPort                   string
	AuthPort                    string
	ProductPort                 string
	CustomerPort                string
	SalePort                    string
	BaseUrl                     string
	DatabaseUsername            string
	DatabasePassword            string
	DatabaseHost                string
	DatabasePort                string
	DatabaseName                string
	EnableDatabaseAutomigration bool
	CloudinaryFolder            string
	CloudinaryCloudName         string
	CloudinaryAPIKey            string
	CLoudinaryAPISecret         string
	SmtpHost                    string
	SmtpUsername                string
	SmtpPassword                string
	SmtpPort                    int
}

func LoadConfig() (config *Config) {
	indexPort := os.Getenv("INDEX_PORT")
	authPort := os.Getenv("AUTH_PORT")
	productPort := os.Getenv("PRODUCT_PORT")
	customerPort := os.Getenv("CUSTOMER_PORT")
	salePort := os.Getenv("SALE_PORT")
	baseUrl := os.Getenv("BASE_URL")
	databaseUsername := os.Getenv("DATABASE_USERNAME")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseHost := os.Getenv("DATABASE_HOST")
	databasePort := os.Getenv("DATABASE_PORT")
	databaseName := os.Getenv("DATABASE_NAME")
	enableDatabaseAutomigration, _ := strconv.ParseBool(os.Getenv("ENABLE_DATABASE_AUTOMIGRATION"))
	cloudinaryFolder := os.Getenv("CLOUDINARY_FOLDER")
	cloudinaryCloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	cloudinaryAPIKey := os.Getenv("CLOUDINARY_API_KEY")
	cLoudinaryAPISecret := os.Getenv("CLOUDINARY_API_SECRET")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	return &Config{
		DatabaseUsername:            databaseUsername,
		DatabasePassword:            databasePassword,
		DatabaseHost:                databaseHost,
		DatabasePort:                databasePort,
		DatabaseName:                databaseName,
		EnableDatabaseAutomigration: enableDatabaseAutomigration,
		BaseUrl:                     baseUrl,
		CloudinaryFolder:            cloudinaryFolder,
		CloudinaryCloudName:         cloudinaryCloudName,
		CloudinaryAPIKey:            cloudinaryAPIKey,
		CLoudinaryAPISecret:         cLoudinaryAPISecret,
		IndexPort:                   indexPort,
		AuthPort:                    authPort,
		ProductPort:                 productPort,
		CustomerPort:                customerPort,
		SalePort:                    salePort,
		SmtpHost:                    smtpHost,
		SmtpUsername:                smtpUsername,
		SmtpPassword:                smtpPassword,
		SmtpPort:                    smtpPort,
	}
}
