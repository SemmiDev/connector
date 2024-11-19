package main

import (
	"errors"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	gl "lab.garudacyber.co.id/g-learning-connector"
	"net/http"

	"strings"
)

type ApiResponse[T any] struct {
	Code     int          `json:"code"`
	Status   string       `json:"status"`
	Message  string       `json:"message"`
	Success  bool         `json:"success"`
	Data     T            `json:"data"`
	Errors   *string      `json:"errors,omitempty"`
	PageInfo *gl.PageInfo `json:"page_info,omitempty"`
}

type ListDataApiResponseWrapper[T any] struct {
	List     []T          `json:"list"`
	PageInfo *gl.PageInfo `json:"page_info,omitempty"`
}

type ErrorHandler interface {
	HTTPStatusCode() int
	Info() string
	Error() string
}

var ErrorTranslator ut.Translator

type ValidationErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func prettyValidationErrorMessage(errs []ValidationErrorMessage) *string {
	builder := strings.Builder{}

	for _, e := range errs {
		builder.WriteString(e.Message)
		builder.WriteString(", ")
	}

	err := strings.TrimSuffix(builder.String(), ", ")
	return &err
}

func HandleError(c *fiber.Ctx, err error) error {
	// Handle if error is a validation error
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		errorMessages := make([]ValidationErrorMessage, 0, len(validationErrors))
		for _, e := range validationErrors {
			field := e.Field()
			message := e.Translate(ErrorTranslator)

			errorMessage := ValidationErrorMessage{Field: field, Message: message}
			errorMessages = append(errorMessages, errorMessage)
		}

		err := prettyValidationErrorMessage(errorMessages)

		return c.Status(http.StatusBadRequest).JSON(ApiResponse[struct{}]{
			Code:    http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Message: *err,
			Success: false,
			Data:    struct{}{},
			Errors:  err,
		})
	}

	// Handle if error is a custom error
	var e ErrorHandler
	if !errors.As(err, &e) {
		errStr := err.Error()
		return c.Status(http.StatusInternalServerError).JSON(ApiResponse[struct{}]{
			Code:    http.StatusInternalServerError,
			Status:  http.StatusText(http.StatusInternalServerError),
			Message: "Terjadi kesalahan pada server",
			Success: false,
			Data:    struct{}{},
			Errors:  &errStr,
		})
	}

	// Otherwise, handle the error
	errStr := e.Error()
	return c.Status(e.HTTPStatusCode()).JSON(ApiResponse[struct{}]{
		Code:    e.HTTPStatusCode(),
		Status:  http.StatusText(e.HTTPStatusCode()),
		Message: e.Info(),
		Success: false,
		Data:    struct{}{},
		Errors:  &errStr,
	})
}
