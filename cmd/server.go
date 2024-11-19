package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"gorm.io/gorm"
	gl "lab.garudacyber.co.id/g-learning-connector"
)

type ApplicationServer struct {
	config *gl.Config
	db     *gorm.DB
	router *fiber.App
}

func NewApplicationServer(db *gorm.DB, config *gl.Config, router *fiber.App) *ApplicationServer {
	app := ApplicationServer{
		config: config,
		db:     db,
		router: router,
	}

	return &app
}

func (a *ApplicationServer) SetupHealthCheckRoutes() {
	a.router.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		ReadinessEndpoint: "/ready",
	}))
}

func (a *ApplicationServer) SetupRoutes() {
	a.router.Get("/api/misca/semesters", a.WithApiKey(), a.ListSemesters)
	a.router.Get("/api/misca/semesters/active", a.WithApiKey(), a.GetActiveSemester)

	a.router.Get("/api/misca/students", a.WithApiKey(), a.ListStudents)
	a.router.Get("/api/misca/students/total", a.WithApiKey(), a.GetTotalStudents)

	a.router.Get("/api/misca/lecturers", a.WithApiKey(), a.ListLecturer)
	a.router.Get("/api/misca/lecturers/total", a.WithApiKey(), a.GetTotalLecturer)
}

func (a *ApplicationServer) Run() {
	host := "0.0.0.0"
	port := fmt.Sprintf("%s", a.config.AppPort)
	hostPort := fmt.Sprintf("%s:%s", host, port)

	err := a.router.Listen(hostPort)
	gl.PanicIfNeeded(err)
}
