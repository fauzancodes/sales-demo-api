package middlewares

import (
	"html/template"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Secure Middleware
func Secure() echo.MiddlewareFunc {
	return middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            0,
		ContentSecurityPolicy: "",
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Path(), "/docs")
		},
	})
}

func StripHTMLMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		for key, values := range c.QueryParams() {
			for i, value := range values {
				sanitizedValue := template.HTMLEscapeString(value)
				sanitizedValue = strings.ReplaceAll(sanitizedValue, "=", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, "<", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, ">", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, "*", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, " AND ", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, " OR ", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, " and ", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, " or ", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, " || ", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, " && ", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, "'", "")
				sanitizedValue = strings.ReplaceAll(sanitizedValue, "&#39;", "")
				values[i] = strip.StripTags(sanitizedValue)
			}
			c.QueryParams()[key] = values
		}

		return next(c)
	}
}
