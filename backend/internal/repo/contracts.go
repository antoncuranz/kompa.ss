// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"backplate/internal/entity"
	"context"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	UserRepo interface {
		GetUsers(ctx context.Context) ([]entity.User, error)
		GetUserByID(ctx context.Context, id int32) (entity.User, error)
	}

	TripsRepo interface {
		GetTrips(ctx context.Context) ([]entity.Trip, error)
		GetTripByID(ctx context.Context, id int32) (entity.Trip, error)
	}

	FlightsRepo interface {
		GetFlights(ctx context.Context) ([]entity.Flight, error)
		GetFlightByID(ctx context.Context, id int32) (entity.Flight, error)
	}

	ActivitiesRepo interface {
		GetActivities(ctx context.Context) ([]entity.Activity, error)
		GetActivityByID(ctx context.Context, id int32) (entity.Activity, error)
	}

	AccommodationRepo interface {
		GetAllAccommodation(ctx context.Context) ([]entity.Accommodation, error)
		GetAccommodationByID(ctx context.Context, id int32) (entity.Accommodation, error)
	}
)
