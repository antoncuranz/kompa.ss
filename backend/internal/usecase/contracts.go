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
		GetFlights(ctx context.Context, tripID int32) ([]entity.Flight, error)
		GetFlightByID(ctx context.Context, tripID int32, id int32) (entity.Flight, error)
		CreateFlight(ctx context.Context, tripID int32, flight request.Flight) (entity.Flight, error)
		UpdateFlight(ctx context.Context, tripID int32, flightID int32, flight request.Flight) error
		DeleteFlight(ctx context.Context, tripID int32, flightID int32) error
	}

	Activities interface {
		GetActivities(ctx context.Context, tripID int32) ([]entity.Activity, error)
		GetActivityByID(ctx context.Context, tripID int32, activityID int32) (entity.Activity, error)
		CreateActivity(ctx context.Context, tripID int32, activity request.Activity) (entity.Activity, error)
		UpdateActivity(ctx context.Context, tripID int32, activityID int32, activity request.Activity) error
		DeleteActivity(ctx context.Context, tripID int32, activityID int32) error
	}

	Accommodation interface {
		GetAllAccommodation(ctx context.Context, tripID int32) ([]entity.Accommodation, error)
		GetAccommodationByID(ctx context.Context, tripID int32, id int32) (entity.Accommodation, error)
		CreateAccommodation(ctx context.Context, tripID int32, accommodation request.Accommodation) (entity.Accommodation, error)
		UpdateAccommodation(ctx context.Context, tripID int32, accommodationID int32, accommodation request.Accommodation) error
		DeleteAccommodation(ctx context.Context, tripID int32, accommodationID int32) error
	}
)
