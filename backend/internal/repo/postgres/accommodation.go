package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"kompass/internal/entity"
	"kompass/pkg/postgres"
	"kompass/pkg/sqlc"
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
		return nil, fmt.Errorf("get all accommodation [tripID=%d]: %w", tripID, err)
	}

	return mapAllAccommodation(rows), nil
}

func (r *AccommodationRepo) GetAccommodationByID(ctx context.Context, tripID int32, accommodationID int32) (entity.Accommodation, error) {
	row, err := r.Queries.GetAccommodationByID(ctx, sqlc.GetAccommodationByIDParams{TripID: tripID, ID: accommodationID})
	if err != nil {
		return entity.Accommodation{}, fmt.Errorf("get accommodation [accommodationID=%d]: %w", accommodationID, err)
	}

	return mapAccommodation(row.Accommodation, mapLocationLeftJoin(row.ID, row.Latitude, row.Longitude)), nil
}

func (r *AccommodationRepo) CreateAccommodation(ctx context.Context, accommodation entity.Accommodation) (entity.Accommodation, error) {
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		return entity.Accommodation{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := r.Queries.WithTx(tx)

	var locationId *int32
	if accommodation.Location != nil {
		persistedLocationId, err := SaveLocation(ctx, qtx, *accommodation.Location)
		if err != nil {
			return entity.Accommodation{}, fmt.Errorf("save location: %w", err)
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
		return entity.Accommodation{}, fmt.Errorf("insert accommodation: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return entity.Accommodation{}, fmt.Errorf("commit tx: %w", err)
	}

	return r.GetAccommodationByID(ctx, accommodation.TripID, accommodationID)
}

func (r *AccommodationRepo) UpdateAccommodation(ctx context.Context, accommodation entity.Accommodation) error {
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := r.Queries.WithTx(tx)

	existingLocationId, err := GetLocationIDOrNilByAccommodationID(ctx, qtx, accommodation.ID)
	if err != nil {
		return fmt.Errorf("get existing location id for accommodation [id=%d]: %w", accommodation.ID, err)
	}

	locationId, err := UpsertOrDeleteLocation(ctx, qtx, existingLocationId, accommodation.Location)
	if err != nil {
		return fmt.Errorf("upsert location [existing=%d]: %w", existingLocationId, err)
	}

	err = qtx.UpdateAccommodation(ctx, sqlc.UpdateAccommodationParams{
		ID:            accommodation.ID,
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
		return fmt.Errorf("update accommodation [id=%d]: %w", accommodation.ID, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (r *AccommodationRepo) DeleteAccommodation(ctx context.Context, tripID int32, accommodationID int32) error {
	return r.Queries.DeleteAccommodationByID(ctx, sqlc.DeleteAccommodationByIDParams{TripID: tripID, ID: accommodationID})
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
		Address:       accommodation.Address,
		Description:   accommodation.Description,
		Price:         accommodation.Price,
		Location:      location,
	}
}
