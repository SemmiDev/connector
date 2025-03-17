package main

import (
	gl "lab.garudacyber.co.id/g-learning-connector"
	"log/slog"
	"os"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// setup logger
	loggerOptions := slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}
	handler := slog.NewJSONHandler(os.Stdout, &loggerOptions)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	config, err := gl.NewConfig()
	gl.PanicIfNeeded(err)

	slog.Info("Config loaded successfully")

	db, err := gl.NewMySQLDatabase(config)
	gl.PanicIfNeeded(err)

	slog.Info("Database connected successfully")

	router := fiber.New(fiber.Config{
		AppName:      config.AppName,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: NewFiberErrorHandler(),
	})

	app := NewApplicationServer(db, logger, config, router)
	app.SetupCommonMiddlewares()
	app.SetupHealthCheckRoutes()
	app.SetupRoutes()

	app.Run()
}
