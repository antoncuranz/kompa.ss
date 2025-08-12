// Package app configures and runs application.
package app

import (
	"fmt"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"os"
	"os/signal"
	"syscall"
	"travel-planner/internal/repo/webapi"
	"travel-planner/internal/usecase"
	"travel-planner/internal/usecase/accommodation"
	"travel-planner/internal/usecase/activities"
	"travel-planner/internal/usecase/attachments"
	"travel-planner/internal/usecase/flights"
	"travel-planner/internal/usecase/trips"
	"travel-planner/internal/usecase/users"
	"travel-planner/pkg/postgres"

	"travel-planner/config"
	"travel-planner/internal/controller/http"
	persistent "travel-planner/internal/repo/postgres"
	"travel-planner/pkg/httpserver"
	"travel-planner/pkg/logger"
)

// Run creates objects via constructors.
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
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
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
	usersUseCase := users.New(persistent.NewUserRepo(pg))
	tripsUseCase := trips.New(persistent.NewTripsRepo(pg))
	flightsUseCase := flights.New(persistent.NewFlightsRepo(pg), webapi.New(cfg.WebApi.AerodataboxApiKey))
	activitiesUseCase := activities.New(persistent.NewActivitiesRepo(pg))
	accommodationUseCase := accommodation.New(persistent.NewAccommodationRepo(pg))
	attachmentsUseCase := attachments.New(persistent.NewAttachmentsRepo(pg))

	return usecase.UseCases{
		Users:         usersUseCase,
		Trips:         tripsUseCase,
		Flights:       flightsUseCase,
		Activities:    activitiesUseCase,
		Accommodation: accommodationUseCase,
		Attachments:   attachmentsUseCase,
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
