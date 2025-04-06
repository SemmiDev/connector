package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// example -> Authorization: Bearer jwtTokenXXX
const (
	authorizationHeaderKey = "Authorization"
	applicationSource      = "glearning"
)

type ApiToken struct {
	ID     int64  `json:"id" gorm:"column:id"`
	Name   string `json:"name" gorm:"column:name"`
	Token  string `json:"token" gorm:"column:token"`
	Active uint8  `json:"active" gorm:"column:active"`
}

func (t *ApiToken) IsActive() bool {
	return t.Active == 1
}

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

		var apiToken ApiToken

		err := a.db.Table("api_key_list").Where("api_key = ?", apiKey).First(&apiToken).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{
					"code":    http.StatusUnauthorized,
					"status":  "Unauthorized",
					"success": false,
					"message": "Api key tidak ditemukan",
				})
			}

			return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"code":    http.StatusInternalServerError,
				"status":  "Internal Server Error",
				"success": false,
				"message": err.Error(),
			})
		}

		return ctx.Next()
	}
}
