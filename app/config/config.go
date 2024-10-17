package config

import (
	"os"
	"strconv"
)

type Config struct {
	SecretKey                     string
	IndexPort                     string
	AuthPort                      string
	ProductPort                   string
	CustomerPort                  string
	SalePort                      string
	PaymentGatewayPort            string
	BaseUrl                       string
	DatabaseUsername              string
	DatabasePassword              string
	DatabaseHost                  string
	DatabasePort                  string
	DatabaseName                  string
	EnableDatabaseAutomigration   bool
	CloudinaryFolder              string
	CloudinaryCloudName           string
	CloudinaryAPIKey              string
	CLoudinaryAPISecret           string
	SmtpHost                      string
	SmtpUsername                  string
	SmtpPassword                  string
	SmtpPort                      int
	MidtransEnv                   string
	MidtransMerchantID            string
	MidtransClientKey             string
	MidtransServerKey             string
	Env                           string
	CustomerImportTemplate        string
	ProductCategoryImportTemplate string
	ProductImportTemplate         string
}

func LoadConfig() (config *Config) {
	secretKey := os.Getenv("SECRET_KEY")
	indexPort := os.Getenv("INDEX_PORT")
	authPort := os.Getenv("AUTH_PORT")
	productPort := os.Getenv("PRODUCT_PORT")
	customerPort := os.Getenv("CUSTOMER_PORT")
	salePort := os.Getenv("SALE_PORT")
	paymentGatewayPort := os.Getenv("PAYMENT_GATEWAY_PORT")
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
	midtransEnv := os.Getenv("MIDTRANS_ENV")
	midtransMerchantID := os.Getenv("MIDTRANS_MERCHANT_ID")
	midtransClientKey := os.Getenv("MIDTRANS_CLIENT_KEY")
	midtransServerKey := os.Getenv("MIDTRANS_SERVER_KEY")
	env := os.Getenv("ENV")
	customerImportTemplate := os.Getenv("CUSTOMER_IMPORT_TEMPLATE")
	productCategoryImportTemplate := os.Getenv("PRODUCT_CATEGORY_IMPORT_TEMPLATE")
	productImportTemplate := os.Getenv("PRODUCT_IMPORT_TEMPLATE")

	return &Config{
		SecretKey:                     secretKey,
		DatabaseUsername:              databaseUsername,
		DatabasePassword:              databasePassword,
		DatabaseHost:                  databaseHost,
		DatabasePort:                  databasePort,
		DatabaseName:                  databaseName,
		EnableDatabaseAutomigration:   enableDatabaseAutomigration,
		BaseUrl:                       baseUrl,
		CloudinaryFolder:              cloudinaryFolder,
		CloudinaryCloudName:           cloudinaryCloudName,
		CloudinaryAPIKey:              cloudinaryAPIKey,
		CLoudinaryAPISecret:           cLoudinaryAPISecret,
		IndexPort:                     indexPort,
		AuthPort:                      authPort,
		ProductPort:                   productPort,
		CustomerPort:                  customerPort,
		SalePort:                      salePort,
		PaymentGatewayPort:            paymentGatewayPort,
		SmtpHost:                      smtpHost,
		SmtpUsername:                  smtpUsername,
		SmtpPassword:                  smtpPassword,
		SmtpPort:                      smtpPort,
		MidtransEnv:                   midtransEnv,
		MidtransMerchantID:            midtransMerchantID,
		MidtransClientKey:             midtransClientKey,
		MidtransServerKey:             midtransServerKey,
		Env:                           env,
		CustomerImportTemplate:        customerImportTemplate,
		ProductCategoryImportTemplate: productCategoryImportTemplate,
		ProductImportTemplate:         productImportTemplate,
	}
}
