package config

import (
	"os"
	"strconv"
)

type Config struct {
	SecretKey                   string
	Port                        string
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
	MidtransEnv                 string
	MidtransMerchantID          string
	MidtransClientKey           string
	MidtransServerKey           string
	IPaymuVA                    string
	IPaymuApiKey                string
	IPaymuBaseURL               string
	XenditBusinessID            string
	XenditSecretKey             string
	XenditPublicKey             string
	XenditWebhookToken          string
	EnableAPIKey                bool
	SpecialApiKey               string
	HMACKey                     string
	BackblazeFolder             string
	BackblazeBucketName         string
	BackblazeKeyID              string
	BackblazeApplicationKey     string
}

func LoadConfig() (config *Config) {
	secretKey := os.Getenv("SECRET_KEY")
	port := os.Getenv("PORT")
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
	ipaymuVA := os.Getenv("IPAYMU_VA")
	ipaymuApiKey := os.Getenv("IPAYMU_API_KEY")
	ipaymuBaseURL := os.Getenv("IPAYMU_BASE_URL")
	xenditBusinessID := os.Getenv("XENDIT_BUSINESS_ID")
	xenditSecretKey := os.Getenv("XENDIT_SECRET_KEY")
	xenditPublicKey := os.Getenv("XENDIT_PUBLIC_KEY")
	xenditWebhookToken := os.Getenv("XENDIT_WEBHOOK_TOKEN")
	enableApiKey, _ := strconv.ParseBool(os.Getenv("ENABLE_API_KEY"))
	specialApiKey := os.Getenv("SPECIAL_API_KEY")
	hmacKey := os.Getenv("HMAC_KEY")
	backblazeFolder := os.Getenv("BACKBLAZE_FOLDER")
	backblazeBucketName := os.Getenv("BACKBLAZE_BUCKET_NAME")
	backblazeKeyID := os.Getenv("BACKBLAZE_KEY_ID")
	backblazeApplicationKey := os.Getenv("BACKBLAZE_APPLICATION_KEY")

	return &Config{
		SecretKey:                   secretKey,
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
		Port:                        port,
		SmtpHost:                    smtpHost,
		SmtpUsername:                smtpUsername,
		SmtpPassword:                smtpPassword,
		SmtpPort:                    smtpPort,
		MidtransEnv:                 midtransEnv,
		MidtransMerchantID:          midtransMerchantID,
		MidtransClientKey:           midtransClientKey,
		MidtransServerKey:           midtransServerKey,
		IPaymuVA:                    ipaymuVA,
		IPaymuApiKey:                ipaymuApiKey,
		IPaymuBaseURL:               ipaymuBaseURL,
		XenditSecretKey:             xenditSecretKey,
		XenditPublicKey:             xenditPublicKey,
		XenditWebhookToken:          xenditWebhookToken,
		XenditBusinessID:            xenditBusinessID,
		EnableAPIKey:                enableApiKey,
		SpecialApiKey:               specialApiKey,
		HMACKey:                     hmacKey,
		BackblazeFolder:             backblazeFolder,
		BackblazeBucketName:         backblazeBucketName,
		BackblazeKeyID:              backblazeKeyID,
		BackblazeApplicationKey:     backblazeApplicationKey,
	}
}
