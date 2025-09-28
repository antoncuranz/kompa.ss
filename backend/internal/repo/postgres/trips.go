package postgres

import (
	"context"
	"fmt"
	"kompass/internal/entity"
	"kompass/pkg/postgres"
	"kompass/pkg/sqlc"
)

type TripsRepo struct {
	*sqlc.Queries
}

func NewTripsRepo(pg *postgres.Postgres) *TripsRepo {
	return &TripsRepo{sqlc.New(pg.Pool)}
}

func (r *TripsRepo) GetTrips(ctx context.Context, userID int32) ([]entity.Trip, error) {
	trips, err := r.Queries.GetTrips(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get trips: %w", err)
	}

	return mapTrips(trips), nil
}

func (r *TripsRepo) GetTripByID(ctx context.Context, userID int32, tripID int32) (entity.Trip, error) {
	trip, err := r.Queries.GetTripByID(ctx, tripID)
	if err != nil {
		return entity.Trip{}, fmt.Errorf("get trip [tripID=%d]: %w", tripID, err)
	}

	return mapTrip(trip), nil
}

func (r *TripsRepo) CreateTrip(ctx context.Context, trip entity.Trip) (entity.Trip, error) {
	tripID, err := r.Queries.InsertTrip(ctx, sqlc.InsertTripParams{
		OwnerID:     trip.OwnerID,
		Name:        trip.Name,
		StartDate:   trip.StartDate,
		EndDate:     trip.EndDate,
		Description: trip.Description,
		ImageUrl:    trip.ImageUrl,
	})
	if err != nil {
		return entity.Trip{}, fmt.Errorf("insert trip: %w", err)
	}

	return r.GetTripByID(ctx, trip.OwnerID, tripID)
}

func (r *TripsRepo) UpdateTrip(ctx context.Context, userID int32, trip entity.Trip) error {
	err := r.Queries.UpdateTrip(ctx, sqlc.UpdateTripParams{
		ID:          trip.ID,
		Name:        trip.Name,
		StartDate:   trip.StartDate,
		EndDate:     trip.EndDate,
		Description: trip.Description,
		ImageUrl:    trip.ImageUrl,
	})
	if err != nil {
		return fmt.Errorf("update trip: %w", err)
	}

	return nil
}

func (r *TripsRepo) DeleteTrip(ctx context.Context, userID int32, tripID int32) error {
	err := r.Queries.DeleteTripByID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("delete trip: %w", err)
	}

	return nil
}

func mapTrips(trips []sqlc.Trip) []entity.Trip {
	result := []entity.Trip{}
	for _, trip := range trips {
		result = append(result, mapTrip(trip))
	}
	return result
}

func mapTrip(trip sqlc.Trip) entity.Trip {
	return entity.Trip{
		ID:          trip.ID,
		OwnerID:     trip.OwnerID,
		Name:        trip.Name,
		StartDate:   trip.StartDate,
		EndDate:     trip.EndDate,
		Description: trip.Description,
		ImageUrl:    trip.ImageUrl,
	}
}
