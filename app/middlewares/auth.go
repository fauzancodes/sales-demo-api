package middlewares

import (
	"net/http"
	"strings"

	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/pkg/jwt"
	"github.com/labstack/echo/v4"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		if token == "" {
			return c.JSON(http.StatusUnauthorized, dto.Response{
				Status:  http.StatusBadRequest,
				Message: "No jwt token provided",
			})
		}

		token = strings.Split(token, " ")[1]
		claims, err := jwt.DecodeToken(token)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, dto.Response{
				Status:  http.StatusUnauthorized,
				Message: "Failed to decode jwt token",
			})
		}

		c.Set("currentUser", claims)
		return next(c)
	}
}
