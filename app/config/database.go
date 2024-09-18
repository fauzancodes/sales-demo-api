package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Database() *gorm.DB {
	host := LoadConfig().DatabaseHost
	user := LoadConfig().DatabaseUsername
	password := LoadConfig().DatabasePassword
	name := LoadConfig().DatabaseName
	port := LoadConfig().DatabasePort

	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, password, name, port)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Print("Failed to connect to database: ", err.Error())
	}

	if LoadConfig().EnableDatabaseAutomigration {
		err = DB.AutoMigrate(
		//database models
		)
		if err != nil {
			log.Print("Failed to migrate database: ", err.Error())
		}
	}

	log.Print("Connected to Database: " + name)

	return DB
}
