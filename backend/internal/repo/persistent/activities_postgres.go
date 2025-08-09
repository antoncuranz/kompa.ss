package persistent

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"travel-planner/internal/entity"
	"travel-planner/pkg/postgres"
	"travel-planner/pkg/sqlc"
)

type ActivitiesRepo struct {
	Db      *pgxpool.Pool
	Queries *sqlc.Queries
}

func NewActivitiesRepo(pg *postgres.Postgres) *ActivitiesRepo {
	return &ActivitiesRepo{
		pg.Pool,
		sqlc.New(pg.Pool),
	}
}

func (r *ActivitiesRepo) GetActivities(ctx context.Context, tripID int32) ([]entity.Activity, error) {
	activities, err := r.Queries.GetActivities(ctx, tripID)
	if err != nil {
		return nil, err
	}

	return mapActivities(activities), nil
}

func (r *ActivitiesRepo) GetActivityByID(ctx context.Context, tripID int32, activityID int32) (entity.Activity, error) {
	row, err := r.Queries.GetActivityByID(ctx, sqlc.GetActivityByIDParams{TripID: tripID, ID: activityID})
	if err != nil {
		return entity.Activity{}, err
	}

	return mapActivity(row.Activity, mapLocationLeftJoin(row.ID, row.Latitude, row.Longitude)), nil
}

func (r *ActivitiesRepo) SaveActivity(ctx context.Context, activity entity.Activity) (entity.Activity, error) {
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		return entity.Activity{}, err
	}
	defer tx.Rollback(ctx)
	qtx := r.Queries.WithTx(tx)

	var locationId *int32
	if activity.Location != nil {
		persistedLocationId, err := SaveLocation(ctx, qtx, *activity.Location)
		if err != nil {
			return entity.Activity{}, err
		}
		locationId = &persistedLocationId
	}

	activityId, err := qtx.InsertActivity(ctx, sqlc.InsertActivityParams{
		TripID:      activity.TripID,
		LocationID:  locationId,
		Name:        activity.Name,
		Date:        activity.Date,
		Time:        activity.Time,
		Address:     activity.Address,
		Description: activity.Description,
		Price:       activity.Price,
	})
	if err != nil {
		return entity.Activity{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return entity.Activity{}, err
	}

	return r.GetActivityByID(ctx, activity.TripID, activityId)
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
		Address:     activity.Address,
		Location:    location,
		Price:       activity.Price,
	}
}
