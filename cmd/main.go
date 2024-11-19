package main

import (
	gl "lab.garudacyber.co.id/g-learning-connector"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config, err := gl.NewConfig()
	gl.PanicIfNeeded(err)

	db, err := gl.NewMySQLDatabase(config)
	gl.PanicIfNeeded(err)

	router := fiber.New(fiber.Config{
		AppName:      config.AppName,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: NewFiberErrorHandler(),
	})

	app := NewApplicationServer(db, config, router)
	app.SetupRoutes()
	app.Run()
}
