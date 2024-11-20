package main

import (
	gl "lab.garudacyber.co.id/g-learning-connector"
	"log/slog"
	"os"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func main() {
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.NewLogLogger(handler, slog.LevelInfo)

	config, err := gl.NewConfig()
	gl.PanicIfNeeded(err)

	logger.Println("Config loaded successfully")

	db, err := gl.NewMySQLDatabase(config)
	gl.PanicIfNeeded(err)

	logger.Println("Database connected successfully")

	router := fiber.New(fiber.Config{
		AppName:      config.AppName,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: NewFiberErrorHandler(),
	})

	app := NewApplicationServer(db, logger, config, router)
	app.SetupHealthCheckRoutes()
	app.SetupRoutes()

	app.Run()
}
