package trips

import (
	"context"
	"kompass/internal/controller/http/v1/request"
	"kompass/internal/entity"
	"kompass/internal/repo"
)

type UseCase struct {
	repo repo.TripsRepo
}

func New(r repo.TripsRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) GetTripByID(ctx context.Context, userID int32, id int32) (entity.Trip, error) {
	return uc.repo.GetTripByID(ctx, userID, id)
}

func (uc *UseCase) GetTrips(ctx context.Context, userID int32) ([]entity.Trip, error) {
	return uc.repo.GetTrips(ctx, userID)
}

func (uc *UseCase) CreateTrip(ctx context.Context, userID int32, trip request.Trip) (entity.Trip, error) {
	return uc.repo.CreateTrip(ctx, entity.Trip{
		OwnerID:     userID,
		Name:        trip.Name,
		StartDate:   trip.StartDate,
		EndDate:     trip.EndDate,
		Description: trip.Description,
		ImageUrl:    trip.ImageUrl,
	})
}

func (uc *UseCase) UpdateTrip(ctx context.Context, userID int32, tripID int32, trip request.Trip) error {
	// TODO: check permissions
	return uc.repo.UpdateTrip(ctx, userID, entity.Trip{
		ID:          tripID,
		Name:        trip.Name,
		StartDate:   trip.StartDate,
		EndDate:     trip.EndDate,
		Description: trip.Description,
		ImageUrl:    trip.ImageUrl,
	})
}

func (uc *UseCase) DeleteTrip(ctx context.Context, userID int32, tripID int32) error {
	return uc.repo.DeleteTrip(ctx, userID, tripID)
}
