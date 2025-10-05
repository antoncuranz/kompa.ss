package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kompass/internal/entity"
	"kompass/pkg/sqlc"
)

func SaveLocation(ctx context.Context, queries *sqlc.Queries, location entity.Location) (int32, error) {
	if location.ID != 0 {
		return location.ID, queries.UpdateLocation(ctx, sqlc.UpdateLocationParams{
			ID:        location.ID,
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
		})
	} else {
		return queries.InsertLocation(ctx, sqlc.InsertLocationParams{
			Latitude:  location.Latitude,
			Longitude: location.Longitude,
		})
	}
}

func GetLocationIDOrNilByActivityID(ctx context.Context, queries *sqlc.Queries, activityID int32) (*int32, error) {
	id, err := queries.GetLocationIDByActivityID(ctx, activityID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &id, err
}

func GetLocationIDOrNilByAccommodationID(ctx context.Context, queries *sqlc.Queries, accommodationID int32) (*int32, error) {
	id, err := queries.GetLocationIDByAccommodationID(ctx, accommodationID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &id, err
}

func UpsertOrDeleteLocation(ctx context.Context, queries *sqlc.Queries, existingLocationId *int32, location *entity.Location) (*int32, error) {
	if location != nil {
		location.ID = *existingLocationId
		savedLocationID, err := SaveLocation(ctx, queries, *location)
		if err != nil {
			return nil, fmt.Errorf("save location: %w", err)
		}
		return &savedLocationID, nil
	} else if existingLocationId != nil {
		if err := DeleteLocation(ctx, queries, *existingLocationId); err != nil {
			return nil, fmt.Errorf("delete location: %w", err)
		}
	}

	return nil, nil
}

func UpdateLocation(ctx context.Context, queries *sqlc.Queries, locationID int32, location entity.Location) error {
	return queries.UpdateLocation(ctx, sqlc.UpdateLocationParams{
		ID:        locationID,
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	})
}

func DeleteLocation(ctx context.Context, queries *sqlc.Queries, locationID int32) error {
	return queries.DeleteLocation(ctx, locationID)
}

func mapLocation(location sqlc.Location) entity.Location {
	return entity.Location{
		ID:        location.ID,
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	}
}

func mapLocationLeftJoin(id *int32, latitude *float32, longitude *float32) *entity.Location {
	if id == nil || latitude == nil || longitude == nil {
		return nil
	}

	return &entity.Location{
		ID:        *id,
		Latitude:  *latitude,
		Longitude: *longitude,
	}
}
