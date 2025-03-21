package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	fiberRecovery "github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
	gl "lab.garudacyber.co.id/g-learning-connector"
	"log/slog"
)

type ApplicationServer struct {
	config *gl.Config
	logger *slog.Logger
	db     *gorm.DB
	router *fiber.App
}

func NewApplicationServer(db *gorm.DB, logger *slog.Logger, config *gl.Config, router *fiber.App) *ApplicationServer {
	app := ApplicationServer{
		config: config,
		logger: logger,
		db:     db,
		router: router,
	}

	return &app
}

func (a *ApplicationServer) SetupCommonMiddlewares() {
	a.router.Use(cors.New())
	a.router.Use(fiberRecovery.New())
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

	a.router.Get("/api/misca/classes", a.WithApiKey(), a.ListKelas)
	a.router.Get("/api/misca/classes/total", a.WithApiKey(), a.TotalKelas)

	a.router.Get("/api/misca/student_classes", a.WithApiKey(), a.ListSimpleStudentKelas)
	a.router.Get("/api/misca/student_classes/total", a.WithApiKey(), a.TotalListSimpleStudentKelas)

	a.router.Get("/api/misca/student_classes_details", a.WithApiKey(), a.ListStudentKelasDetails)
	a.router.Get("/api/misca/student_classes_details/total", a.WithApiKey(), a.GetTotalKelasDetails)
}

func (a *ApplicationServer) Run() {
	host := "0.0.0.0"
	port := fmt.Sprintf("%s", a.config.AppPort)
	hostPort := fmt.Sprintf("%s:%s", host, port)

	a.logger.With(slog.String("host", host), slog.String("port", port)).Info("Server started")

	err := a.router.Listen(hostPort)
	gl.PanicIfNeeded(err)
}
