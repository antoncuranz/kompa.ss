package persistent

import (
	"context"
	"travel-planner/internal/entity"
	"travel-planner/pkg/sqlc"
)

func SaveLocation(ctx context.Context, queries *sqlc.Queries, location entity.Location) (int32, error) {
	return queries.InsertLocation(ctx, sqlc.InsertLocationParams{
		Latitude:  location.Latitude,
		Longitude: location.Longitude,
	})
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
