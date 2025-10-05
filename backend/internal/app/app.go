// Package app configures and runs application.
package app

import (
	"fmt"
	"kompass/internal/controller/http/v1/response"
	"kompass/internal/repo/webapi"
	"kompass/internal/usecase"
	"kompass/internal/usecase/accommodation"
	"kompass/internal/usecase/activities"
	"kompass/internal/usecase/attachments"
	"kompass/internal/usecase/flights"
	"kompass/internal/usecase/geocoding"
	"kompass/internal/usecase/trains"
	"kompass/internal/usecase/transportation"
	"kompass/internal/usecase/trips"
	"kompass/internal/usecase/users"
	"kompass/pkg/postgres"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"kompass/config"
	"kompass/internal/controller/http"
	persistent "kompass/internal/repo/postgres"
	"kompass/pkg/httpserver"
	"kompass/pkg/logger"
)

func Run(cfg *config.Config) {
	log := logger.New(cfg.Log.Level)

	// Repository
	pg, err := connectAndMigrate(cfg.PG)
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - ConnectAndMigrate: %w", err))
	}
	defer pg.Close()

	// Use-Case
	useCases := createUseCases(cfg, pg)

	// HTTP Server
	httpServer := httpserver.New(
		httpserver.Port(cfg.HTTP.Port),
		httpserver.Prefork(cfg.HTTP.UsePreforkMode),
		httpserver.ErrorHandler(response.ErrorHandler),
	)
	http.NewRouter(httpServer.App, cfg, useCases, log)
	httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		log.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}

func createUseCases(cfg *config.Config, pg *postgres.Postgres) usecase.UseCases {
	flightsRepo := persistent.NewFlightsRepo(pg)
	transportationRepo := persistent.NewTransportationRepo(pg, flightsRepo, persistent.NewTrainsRepo(pg))
	ors := webapi.NewOpenRouteServiceWebAPI(cfg.WebApi)

	usersUseCase := users.New(persistent.NewUserRepo(pg))
	tripsUseCase := trips.New(persistent.NewTripsRepo(pg))
	transportationUseCase := transportation.New(transportationRepo, ors)
	flightsUseCase := flights.New(transportationRepo, flightsRepo, webapi.New(cfg.WebApi))
	trainsUseCase := trains.New(transportationRepo, webapi.NewDbVendoWebAPI(cfg.WebApi))
	activitiesUseCase := activities.New(persistent.NewActivitiesRepo(pg), tripsUseCase)
	accommodationUseCase := accommodation.New(persistent.NewAccommodationRepo(pg), tripsUseCase)
	attachmentsUseCase := attachments.New(persistent.NewAttachmentsRepo(pg))
	geocodingUseCase := geocoding.New(trainsUseCase, ors)

	return usecase.UseCases{
		Users:          usersUseCase,
		Geocoding:      geocodingUseCase,
		Trips:          tripsUseCase,
		Transportation: transportationUseCase,
		Flights:        flightsUseCase,
		Trains:         trainsUseCase,
		Activities:     activitiesUseCase,
		Accommodation:  accommodationUseCase,
		Attachments:    attachmentsUseCase,
	}
}

func connectAndMigrate(cfg config.PG) (*postgres.Postgres, error) {
	pg, err := postgres.New(cfg.URL, postgres.MaxPoolSize(cfg.PoolMax))
	if err != nil {
		return nil, err
	}
	if err := goose.SetDialect("postgres"); err != nil {
		return nil, err
	}
	db := stdlib.OpenDBFromPool(pg.Pool)
	if err := goose.Up(db, "migrations"); err != nil {
		return nil, err
	}
	return pg, nil
}
