package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/fauzancodes/sales-demo-api/app/config"
	"github.com/fauzancodes/sales-demo-api/app/dto"
	"github.com/fauzancodes/sales-demo-api/app/models"
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

func CheckAPIKey(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if config.LoadConfig().EnableAPIKey {
			apiKey := c.Request().Header.Get("X-API-KEY")
			if apiKey == "" {
				fmt.Println("Failed to check api key: no api key in header")
				return c.JSON(http.StatusForbidden, dto.Response{
					Status:  http.StatusForbidden,
					Message: "Forbidden",
				})
			}
			if apiKey == config.LoadConfig().SpecialApiKey {
				return next(c)
			}

			secretKey, receivedHMAC, err := DecodeAPIKeyBase64(apiKey)
			if err != nil {
				fmt.Println("Failed to check api key: ", err.Error())
				return c.JSON(http.StatusForbidden, dto.Response{
					Status:  http.StatusForbidden,
					Message: "Forbidden",
				})
			}

			secret := config.LoadConfig().HMACKey

			hmacVerified, expectedHMAC, err := VerifyAPIKeyHMAC(secretKey, receivedHMAC, secret)
			if err != nil {
				fmt.Println("Failed to check api key: ", err.Error())
				return c.JSON(http.StatusForbidden, dto.Response{
					Status:  http.StatusForbidden,
					Message: "Forbidden",
				})
			}
			var usedApiKey models.SDAUsedApiKey
			err = config.DB.Debug().Where("secret_key = ?", secretKey).First(&usedApiKey).Error
			if err != nil {
				fmt.Println("Check api key:", err.Error())
			}
			if usedApiKey.ID > 0 {
				fmt.Println("Failed to check api key: api key already used")
				return c.JSON(http.StatusForbidden, dto.Response{
					Status:  http.StatusForbidden,
					Message: "Forbidden",
				})
			}

			err = config.DB.Create(&models.SDAUsedApiKey{
				SecretKey:    secretKey,
				Base64Key:    apiKey,
				ReceivedHMAC: receivedHMAC,
				ExpectedHMAC: expectedHMAC,
			}).Error
			if err != nil {
				fmt.Println("Failed to save api key: ", err.Error())
			}

			if !hmacVerified {
				fmt.Println("Failed to check api key: failed to verify hmac")
				return c.JSON(http.StatusForbidden, dto.Response{
					Status:  http.StatusForbidden,
					Message: "Forbidden",
				})
			}
		}

		return next(c)
	}
}

func ComputeAPIKeyHMAC(secretKey, secret string) (response string, err error) {
	h := hmac.New(sha256.New, []byte(secret))
	_, err = h.Write([]byte(secretKey))
	response = hex.EncodeToString(h.Sum(nil))

	return
}

func DecodeAPIKeyBase64(encodedKey string) (secretKey string, hmacSignature string, err error) {
	decoded, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return
	}

	parts := strings.Split(string(decoded), ":")
	if len(parts) != 2 {
		err = errors.New("invalid encoded key")
		return
	}

	secretKey = parts[0]
	hmacSignature = parts[1]

	return
}

func VerifyAPIKeyHMAC(secretKey, receivedHMAC, secret string) (response bool, expectedHMAC string, err error) {
	expectedHMAC, err = ComputeAPIKeyHMAC(secretKey, secret)
	if err != nil {
		return
	}
	response = hmac.Equal([]byte(expectedHMAC), []byte(receivedHMAC))

	return
}
