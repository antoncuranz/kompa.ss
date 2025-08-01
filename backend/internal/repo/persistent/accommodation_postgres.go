package persistent

import (
	"context"
	"travel-planner/internal/entity"
	"travel-planner/pkg/postgres"
	"travel-planner/pkg/sqlc"
)

type AccommodationRepo struct {
	*sqlc.Queries
}

func NewAccommodationRepo(pg *postgres.Postgres) *AccommodationRepo {
	return &AccommodationRepo{sqlc.New(pg.Pool)}
}

func (r *AccommodationRepo) GetAllAccommodation(ctx context.Context) ([]entity.Accommodation, error) {
	accommodation, err := r.Queries.GetAllAccommodation(ctx)
	if err != nil {
		return nil, err
	}

	return mapAllAccommodation(accommodation), nil
}

func (r *AccommodationRepo) GetAccommodationByID(ctx context.Context, id int32) (entity.Accommodation, error) {
	accommodation, err := r.Queries.GetAccommodationByID(ctx, id)
	if err != nil {
		return entity.Accommodation{}, err
	}

	return mapAccommodation(accommodation), nil
}

func mapAllAccommodation(accommodation []sqlc.Accommodation) []entity.Accommodation {
	result := []entity.Accommodation{}
	for _, accommodation := range accommodation {
		result = append(result, mapAccommodation(accommodation))
	}
	return result
}

func mapAccommodation(accommodation sqlc.Accommodation) entity.Accommodation {
	return entity.Accommodation{
		ID:            accommodation.ID,
		TripID:        accommodation.TripID,
		Name:          accommodation.Name,
		ArrivalDate:   accommodation.ArrivalDate,
		DepartureDate: accommodation.DepartureDate,
		CheckInTime:   accommodation.CheckInTime,
		CheckOutTime:  accommodation.CheckOutTime,
		Description:   accommodation.Description,
		Location:      accommodation.Location,
		Price:         accommodation.Price,
	}
}
