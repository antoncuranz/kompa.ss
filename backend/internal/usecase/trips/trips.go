package trips

import (
	"cloud.google.com/go/civil"
	"context"
	"fmt"
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

func (uc *UseCase) GetTripByID(ctx context.Context, id int32) (entity.Trip, error) {
	return uc.repo.GetTripByID(ctx, id)
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

func (uc *UseCase) VerifyDatesInBounds(ctx context.Context, tripID int32, dates ...civil.Date) error {
	for _, d := range dates {
		if err := uc.verifyDateInBounds(ctx, tripID, d); err != nil {
			return err
		}
	}

	return nil
}

func (uc *UseCase) verifyDateInBounds(ctx context.Context, tripID int32, date civil.Date) error {
	trip, err := uc.GetTripByID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("get trip [id=%d]: %w", tripID, err)
	}

	if trip.StartDate.After(date) {
		return fmt.Errorf("%s is before start date (%s)", date.String(), trip.StartDate.String())
	} else if trip.EndDate.Before(date) {
		return fmt.Errorf("%s is after end date (%s)", date.String(), trip.EndDate.String())
	}

	return nil
}
