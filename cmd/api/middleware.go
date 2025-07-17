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

		// Pilih tabel berdasarkan URL
		var tableName string
		url := ctx.OriginalURL()

		instansi := "MISCA"

		if strings.Contains(url, "smart") {
			tableName = "setting_app"
			instansi = "SMART"
		} else if strings.Contains(url, "misca") {
			tableName = "setting_pt"
		} else {
			return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
				"code":    http.StatusBadRequest,
				"status":  "Bad Request",
				"success": false,
				"message": "Unknown app source in URL",
			})
		}

		// Query secret dari tabel yang dipilih
		var secret string
		err := a.db.Table(tableName).Where("param = ?", "secret_smartthink").Select("value").Scan(&secret).Error
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

		// save to context
		ctx.Locals("tipe_instansi", instansi)

		return ctx.Next()
	}
}
