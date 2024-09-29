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
	routes.AuthRoute(app)

	port := config.LoadConfig().AuthPort

	log.Printf("Server: " + config.LoadConfig().BaseURL + ":" + port)
	app.Logger.Fatal(app.Start(":" + port))
}
