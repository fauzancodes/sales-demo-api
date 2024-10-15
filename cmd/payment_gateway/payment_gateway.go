package main

import (
	"log"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/routes"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()

	config.Database()
	routes.PaymentGatewayRoute(app)

	port := config.LoadConfig().PaymentGatewayPort

	log.Printf("Server: " + config.LoadConfig().BaseUrl + ":" + port)
	app.Logger.Fatal(app.Start(":" + port))
}
