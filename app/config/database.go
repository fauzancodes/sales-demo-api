package config

import (
	"fmt"
	"log"

	"github.com/fauzancodes/sales-demo-api/app/models"
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
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if LoadConfig().EnableDatabaseAutomigration {
		go RunAutoMigration()
	}

	log.Printf("Connected to Database: %v", name)

	return DB
}

func RunAutoMigration() {
	err := DB.AutoMigrate(
		&models.SDAUser{},
		&models.SDAProductCategory{},
		&models.SDAProduct{},
		&models.SDAProductStock{},
		&models.SDACustomer{},
		&models.SDASale{},
		&models.SDASaleDetail{},
		&models.SDAMidtransPaymentMethod{},
		&models.SDAMidtransSalePayment{},
		&models.SDAIPaymuPaymentMethod{},
		&models.SDAIPaymuSalePayment{},
		&models.SDAXenditPaymentMethod{},
		&models.SDAXenditSalePayment{},
		&models.SDAUsedApiKey{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
