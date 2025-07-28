package persistent

import (
	"backplate/internal/entity"
	"backplate/pkg/postgres"
	"backplate/pkg/sqlc"
	"context"
)

type TripsRepo struct {
	*sqlc.Queries
}

func NewTripsRepo(pg *postgres.Postgres) *TripsRepo {
	return &TripsRepo{sqlc.New(pg.Pool)}
}

func (r *TripsRepo) GetTrips(ctx context.Context) ([]entity.Trip, error) {
	trips, err := r.Queries.GetTrips(ctx)
	if err != nil {
		return nil, err
	}

	return mapTrips(trips), nil
}

func (r *TripsRepo) GetTripByID(ctx context.Context, id int32) (entity.Trip, error) {
	trip, err := r.Queries.GetTripByID(ctx, id)
	if err != nil {
		return entity.Trip{}, err
	}

	return mapTrip(trip), nil
}

func mapTrips(trips []sqlc.Trip) []entity.Trip {
	var result []entity.Trip
	for _, trip := range trips {
		result = append(result, mapTrip(trip))
	}
	return result
}

func mapTrip(trip sqlc.Trip) entity.Trip {
	return entity.Trip{
		ID:          trip.ID,
		Name:        trip.Name,
		Description: trip.Description,
		StartDate:   trip.StartDate,
		EndDate:     trip.EndDate,
	}
}
