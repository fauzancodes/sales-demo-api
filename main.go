package handler //for vercel
// package main

import (
	"log"
	"net/http"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/middlewares"
	"github.com/fauzancodes/sales-demo-api/app/routes"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

func main() {
	app := Init()

	port := config.LoadConfig().Port

	log.Printf("Server: " + config.LoadConfig().BaseUrl + ":" + port)
	app.Logger.Fatal(app.Start(":" + port))
}

func Main(w http.ResponseWriter, r *http.Request) {
	e := Init()

	e.ServeHTTP(w, r)
}

func Init() *echo.Echo {
	app := echo.New()

	app.Use(middlewares.Cors())
	app.Use(middlewares.Gzip())
	app.Use(middlewares.Logger())
	app.Use(middlewares.Secure())
	app.Use(middlewares.Recover())

	config.Database()

	routes.Route(app)

	return app
}
