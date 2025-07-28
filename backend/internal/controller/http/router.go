// Package v1 implements routing paths. Each services in own file.
package http

import (
	"backplate/internal/usecase"
	"net/http"

	"backplate/config"
	_ "backplate/docs" // Swagger docs.
	"backplate/internal/controller/http/middleware"
	v1 "backplate/internal/controller/http/v1"
	"backplate/pkg/logger"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// NewRouter -.
// Swagger spec:
// @title       Backplate API
// @description Using a translation service as an example
// @version     1.0
// @host        127.0.0.1:8080
// @BasePath    /api/v1
func NewRouter(app *fiber.App, cfg *config.Config, useCases usecase.UseCases, log logger.Interface) {
	// Options
	app.Use(middleware.Logger(log))
	app.Use(middleware.Recovery(log))

	// Prometheus metrics
	if cfg.Metrics.Enabled {
		prometheus := fiberprometheus.New("my-service-name")
		prometheus.RegisterAt(app, "/metrics")
		app.Use(prometheus.Middleware)
	}

	// Swagger
	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	// K8s probe
	app.Get("/healthz", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	// Routers
	apiV1Group := app.Group("/api/v1")
	{
		v1.NewUserRoutes(apiV1Group, useCases.Users, log)
		v1.NewTripRoutes(apiV1Group, useCases.Trips, log)
		v1.NewFlightRoutes(apiV1Group, useCases.Flights, log)
		v1.NewActivityRoutes(apiV1Group, useCases.Activities, log)
		v1.NewAccommodationRoutes(apiV1Group, useCases.Accommodation, log)
	}
}
