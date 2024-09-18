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
	config.Redis()
	routes.RouteInit(app)

	log.Printf("Server: " + config.LoadConfig().BaseURL + ":" + config.LoadConfig().Port)
	app.Logger.Fatal(app.Start(":" + config.LoadConfig().Port))
}
