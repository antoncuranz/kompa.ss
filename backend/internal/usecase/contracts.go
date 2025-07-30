// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	"travel-planner/internal/controller/http/v1/request"
	"travel-planner/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	UseCases struct {
		Users         Users
		Trips         Trips
		Flights       Flights
		Activities    Activities
		Accommodation Accommodation
	}

	Users interface {
		GetUsers(ctx context.Context) ([]entity.User, error)
		GetUserByID(ctx context.Context, id int32) (entity.User, error)
	}

	Trips interface {
		GetTrips(ctx context.Context) ([]entity.Trip, error)
		GetTripByID(ctx context.Context, id int32) (entity.Trip, error)
	}

	Flights interface {
		GetFlights(ctx context.Context) ([]entity.Flight, error)
		GetFlightByID(ctx context.Context, id int32) (entity.Flight, error)
		CreateFlight(ctx context.Context, flight request.Flight) (entity.Flight, error)
	}

	Activities interface {
		GetActivities(ctx context.Context) ([]entity.Activity, error)
		GetActivityByID(ctx context.Context, id int32) (entity.Activity, error)
	}

	Accommodation interface {
		GetAllAccommodation(ctx context.Context) ([]entity.Accommodation, error)
		GetAccommodationByID(ctx context.Context, id int32) (entity.Accommodation, error)
	}
)
