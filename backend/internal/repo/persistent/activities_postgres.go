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
	activity, err := r.Queries.GetActivityByID(ctx, id)
	if err != nil {
		return entity.Activity{}, err
	}

	return mapActivity(activity), nil
}

func mapActivities(activities []sqlc.Activity) []entity.Activity {
	result := []entity.Activity{}
	for _, activity := range activities {
		result = append(result, mapActivity(activity))
	}
	return result
}

func mapActivity(activity sqlc.Activity) entity.Activity {
	return entity.Activity{
		ID:          activity.ID,
		TripID:      activity.TripID,
		Name:        activity.Name,
		Description: activity.Description,
		Date:        activity.Date,
	}
}
