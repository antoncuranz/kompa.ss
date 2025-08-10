package trips

import (
	"context"
	"travel-planner/internal/controller/http/v1/request"
	"travel-planner/internal/entity"
	"travel-planner/internal/repo"
)

type UseCase struct {
	repo repo.TripsRepo
}

func New(r repo.TripsRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) GetTripByID(ctx context.Context, id int32) (entity.Trip, error) {
	return uc.repo.GetTripByID(ctx, id)
}

func (uc *UseCase) GetTrips(ctx context.Context) ([]entity.Trip, error) {
	return uc.repo.GetTrips(ctx)
}

func (uc *UseCase) CreateTrip(ctx context.Context, trip request.Trip) (entity.Trip, error) {
	return uc.repo.CreateTrip(ctx, entity.Trip{
		Name:        trip.Name,
		StartDate:   trip.StartDate,
		EndDate:     trip.EndDate,
		Description: trip.Description,
		ImageUrl:    trip.ImageUrl,
	})
}

func (uc *UseCase) UpdateTrip(ctx context.Context, tripID int32, trip request.Trip) error {
	return uc.repo.UpdateTrip(ctx, entity.Trip{
		ID:          tripID,
		Name:        trip.Name,
		StartDate:   trip.StartDate,
		EndDate:     trip.EndDate,
		Description: trip.Description,
		ImageUrl:    trip.ImageUrl,
	})
}

func (uc *UseCase) DeleteTrip(ctx context.Context, tripID int32) error {
	return uc.repo.DeleteTrip(ctx, tripID)
}
