package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"net/http"
)

func NewFiberErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		statusCode := fiber.StatusInternalServerError

		var e *fiber.Error
		if errors.As(err, &e) {
			statusCode = e.Code
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"code":    statusCode,
			"status":  http.StatusText(statusCode),
			"message": err.Error(),
			"success": false,
			"data":    struct{}{},
			"errors":  nil,
		})
	}
}
