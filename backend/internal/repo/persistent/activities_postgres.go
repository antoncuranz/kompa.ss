package persistent

import (
	"context"
	"travel-planner/internal/entity"
	"travel-planner/pkg/postgres"
	"travel-planner/pkg/sqlc"
)

type ActivitiesRepo struct {
	*sqlc.Queries
}

func NewActivitiesRepo(pg *postgres.Postgres) *ActivitiesRepo {
	return &ActivitiesRepo{sqlc.New(pg.Pool)}
}

func (r *ActivitiesRepo) GetActivities(ctx context.Context) ([]entity.Activity, error) {
	activities, err := r.Queries.GetActivities(ctx)
	if err != nil {
		return nil, err
	}

	return mapActivities(activities), nil
}

func (r *ActivitiesRepo) GetActivityByID(ctx context.Context, id int32) (entity.Activity, error) {
	row, err := r.Queries.GetActivityByID(ctx, id)
	if err != nil {
		return entity.Activity{}, err
	}

	return mapActivity(row.Activity, mapLocationLeftJoin(row.ID, row.Latitude, row.Longitude)), nil
}

func mapActivities(rows []sqlc.GetActivitiesRow) []entity.Activity {
	result := []entity.Activity{}
	for _, row := range rows {
		result = append(result, mapActivity(row.Activity, mapLocationLeftJoin(row.ID, row.Latitude, row.Longitude)))
	}
	return result
}

func mapActivity(activity sqlc.Activity, location *entity.Location) entity.Activity {
	return entity.Activity{
		ID:          activity.ID,
		TripID:      activity.TripID,
		Name:        activity.Name,
		Date:        activity.Date,
		Time:        activity.Time,
		Description: activity.Description,
		Location:    location,
	}
}
