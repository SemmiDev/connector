package main

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// example -> Authorization: Bearer jwtTokenXXX
const (
	authorizationHeaderKey = "Authorization"
	applicationSource      = "glearning"

	instansiTypeKey   = "tipe_instansi"
	instansiTypeMisca = "MISCA"
	instansiTypeSmart = "SMART"
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
		var instansi string

		// Coba dulu ke setting_pt
		err := a.db.Table("setting_pt").
			Where("param = ?", "secret_smartthink").
			Select("value").
			Scan(&secret).Error

		if err == nil && secret != "" {
			instansi = instansiTypeMisca
		} else {
			// Kalau tidak ada, fallback ke setting_app
			err = a.db.Table("setting_app").
				Where("param = ?", "secret_smartthink").
				Select("value").
				Scan(&secret).Error
			if err != nil {
				return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"code":    http.StatusInternalServerError,
					"status":  "Internal Server Error",
					"success": false,
					"message": err.Error(),
				})
			}
			instansi = instansiTypeSmart
		}

		if secret == "" {
			return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"code":    http.StatusUnauthorized,
				"status":  "Unauthorized",
				"success": false,
				"message": "Secret tidak ditemukan",
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

		ctx.Locals(instansiTypeKey, instansi)
		return ctx.Next()
	}
}

func (a *ApplicationServer) SetupCommonMiddlewares() {
	a.router.Use(cors.New())
	a.router.Use(recover.New())
}

func IsSmartInstansi(c *fiber.Ctx) bool {
	tipe := c.Locals(instansiTypeKey)
	return tipe == instansiTypeSmart
}
