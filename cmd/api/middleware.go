package main

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// example -> Authorization: Bearer jwtTokenXXX
const (
	authorizationHeaderKey = "Authorization"
	applicationSource      = "glearning"
)

func (a *ApplicationServer) WithApiKey() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authorizationHeader := strings.TrimSpace(ctx.Get(authorizationHeaderKey))
		if authorizationHeader == "" {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"code":    http.StatusUnauthorized,
				"status":  "Unauthorized",
				"success": false,
				"message": "Api key tidak ditemukan",
			})
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"code":    http.StatusUnauthorized,
				"status":  "Unauthorized",
				"success": false,
				"message": "Format api key salah",
			})
		}

		apiKey := fields[1]

		var secret string
		err := a.db.Table("setting_pt").Where("param = ?", "secret_smartthink").Select("value").Scan(&secret).Error
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"status":  "Internal Server Error",
				"success": false,
				"message": err.Error(),
			})
		}

		if secret != apiKey {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"code":    http.StatusUnauthorized,
				"status":  "Unauthorized",
				"success": false,
				"message": "Api key tidak sesuai",
			})
		}

		return ctx.Next()
	}
}
