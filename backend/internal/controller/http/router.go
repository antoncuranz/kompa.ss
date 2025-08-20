// Package v1 implements routing paths. Each services in own file.
package http

import (
	"kompass/internal/usecase"
	"net/http"

	fiberSwagger "github.com/swaggo/fiber-swagger"

	"kompass/config"
	_ "kompass/docs" // Swagger docs.
	"kompass/internal/controller/http/middleware"
	v1 "kompass/internal/controller/http/v1"
	"kompass/pkg/logger"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
)

// NewRouter -.
// Swagger spec:
// @title       Kompa.ss API
// @description Using a translation service as an example
// @version     1.0
// @servers.url http://127.0.0.1:8080/api/v1
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
		app.Get("/swagger/*", fiberSwagger.WrapHandler)
	}

	// K8s probe
	app.Get("/healthz", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	// Routers
	apiV1Group := app.Group("/api/v1")
	{
		v1.NewUserRoutes(apiV1Group, useCases.Users, log)

		tripsV1Group := v1.NewTripRoutes(apiV1Group, useCases.Trips, log)
		{
			v1.NewFlightRoutes(tripsV1Group, useCases.Flights, log)
			v1.NewActivityRoutes(tripsV1Group, useCases.Activities, log)
			v1.NewAccommodationRoutes(tripsV1Group, useCases.Accommodation, log)
			v1.NewAttachmentRoutes(tripsV1Group, useCases.Attachments, log)
		}
	}
}
