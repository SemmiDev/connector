package main

import (
	"log/slog"
	"os"

	gl "lab.garudacyber.co.id/g-learning-connector"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// load config
	config, err := gl.NewConfig()
	gl.PanicIfNeeded(err)

	// set up logging
	loggerOptions := slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}
	handler := slog.NewJSONHandler(os.Stdout, &loggerOptions)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("Config loaded successfully")

	// set up database
	db, err := gl.NewMySQLDatabase(config)
	gl.PanicIfNeeded(err)

	slog.Info("Database connected successfully")

	// set up route
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
