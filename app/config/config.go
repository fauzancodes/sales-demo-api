package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port                        string
	BaseURL                     string
	CacheURL                    string
	CachePassword               string
	LoggerLevel                 string
	ContextTimeout              int
	DatabaseUsername            string
	DatabasePassword            string
	DatabaseHost                string
	DatabasePort                string
	DatabaseName                string
	EnableDatabaseAutomigration bool
}

func LoadConfig() (config *Config) {
	cacheURL := os.Getenv("CACHE_URL")
	cachePassword := os.Getenv("CACHE_PASSWORD")
	loggerLevel := os.Getenv("LOGGER_LEVEL")
	contextTimeout, _ := strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT"))
	port := os.Getenv("PORT")
	baseUrl := os.Getenv("BASE_URL")
	databaseUsername := os.Getenv("DATABASE_USERNAME")
	databasePassword := os.Getenv("DATABASE_PASSWORD")
	databaseHost := os.Getenv("DATABASE_HOST")
	databasePort := os.Getenv("DATABASE_PORT")
	databaseName := os.Getenv("DATABASE_NAME")
	enableDatabaseAutomigration, _ := strconv.ParseBool(os.Getenv("ENABLE_DATABASE_AUTOMIGRATION"))

	return &Config{
		CacheURL:                    cacheURL,
		CachePassword:               cachePassword,
		Port:                        port,
		LoggerLevel:                 loggerLevel,
		ContextTimeout:              contextTimeout,
		DatabaseUsername:            databaseUsername,
		DatabasePassword:            databasePassword,
		DatabaseHost:                databaseHost,
		DatabasePort:                databasePort,
		DatabaseName:                databaseName,
		EnableDatabaseAutomigration: enableDatabaseAutomigration,
		BaseURL:                     baseUrl,
	}
}
