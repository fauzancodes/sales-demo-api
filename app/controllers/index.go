package controllers

import (
	"github.com/labstack/echo/v4"
)

func Index(c echo.Context) error {
	return c.HTML(200, `
		<h1>Sales Demo API</h1>
		<a href="/assets/Sale Demo API.postman_collection.json" download>
			Click here to download the postman collection!
		</a>
	`)
}
