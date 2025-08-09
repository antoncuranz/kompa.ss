package persistent

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"travel-planner/internal/entity"
	"travel-planner/pkg/postgres"
	"travel-planner/pkg/sqlc"
)

type AccommodationRepo struct {
	Db      *pgxpool.Pool
	Queries *sqlc.Queries
}

func NewAccommodationRepo(pg *postgres.Postgres) *AccommodationRepo {
	return &AccommodationRepo{
		pg.Pool,
		sqlc.New(pg.Pool),
	}
}

func (r *AccommodationRepo) GetAllAccommodation(ctx context.Context, tripID int32) ([]entity.Accommodation, error) {
	rows, err := r.Queries.GetAllAccommodation(ctx, tripID)
	if err != nil {
		return nil, err
	}

	return mapAllAccommodation(rows), nil
}

func (r *AccommodationRepo) GetAccommodationByID(ctx context.Context, tripID int32, accommodationID int32) (entity.Accommodation, error) {
	row, err := r.Queries.GetAccommodationByID(ctx, sqlc.GetAccommodationByIDParams{TripID: tripID, ID: accommodationID})
	if err != nil {
		return entity.Accommodation{}, err
	}

	return mapAccommodation(row.Accommodation, mapLocationLeftJoin(row.ID, row.Latitude, row.Longitude)), nil
}

func (r *AccommodationRepo) SaveAccommodation(ctx context.Context, accommodation entity.Accommodation) (entity.Accommodation, error) {
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		return entity.Accommodation{}, err
	}
	defer tx.Rollback(ctx)
	qtx := r.Queries.WithTx(tx)

	var locationId *int32
	if accommodation.Location != nil {
		persistedLocationId, err := SaveLocation(ctx, qtx, *accommodation.Location)
		if err != nil {
			return entity.Accommodation{}, err
		}
		locationId = &persistedLocationId
	}

	accommodationID, err := qtx.InsertAccommodation(ctx, sqlc.InsertAccommodationParams{
		TripID:        accommodation.TripID,
		LocationID:    locationId,
		Name:          accommodation.Name,
		ArrivalDate:   accommodation.ArrivalDate,
		DepartureDate: accommodation.DepartureDate,
		CheckInTime:   accommodation.CheckInTime,
		CheckOutTime:  accommodation.CheckOutTime,
		Address:       accommodation.Address,
		Description:   accommodation.Description,
		Price:         accommodation.Price,
	})
	if err != nil {
		return entity.Accommodation{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return entity.Accommodation{}, err
	}

	return r.GetAccommodationByID(ctx, accommodation.ID, accommodationID)
}

func mapAllAccommodation(accommodation []sqlc.GetAllAccommodationRow) []entity.Accommodation {
	result := []entity.Accommodation{}
	for _, row := range accommodation {
		result = append(result, mapAccommodation(row.Accommodation, mapLocationLeftJoin(row.ID, row.Latitude, row.Longitude)))
	}
	return result
}

func mapAccommodation(accommodation sqlc.Accommodation, location *entity.Location) entity.Accommodation {
	return entity.Accommodation{
		ID:            accommodation.ID,
		TripID:        accommodation.TripID,
		Name:          accommodation.Name,
		ArrivalDate:   accommodation.ArrivalDate,
		DepartureDate: accommodation.DepartureDate,
		CheckInTime:   accommodation.CheckInTime,
		CheckOutTime:  accommodation.CheckOutTime,
		Description:   accommodation.Description,
		Price:         accommodation.Price,
		Location:      location,
	}
}
