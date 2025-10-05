package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"kompass/internal/entity"
	"kompass/pkg/postgres"
	"kompass/pkg/sqlc"
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
		return nil, fmt.Errorf("get all activities [tripID=%d]: %w", tripID, err)
	}

	return mapActivities(activities), nil
}

func (r *ActivitiesRepo) GetActivityByID(ctx context.Context, tripID int32, activityID int32) (entity.Activity, error) {
	row, err := r.Queries.GetActivityByID(ctx, sqlc.GetActivityByIDParams{TripID: tripID, ID: activityID})
	if err != nil {
		return entity.Activity{}, fmt.Errorf("get activity [id=%d]: %w", activityID, err)
	}

	return mapActivity(row.Activity, mapLocationLeftJoin(row.ID, row.Latitude, row.Longitude)), nil
}

func (r *ActivitiesRepo) CreateActivity(ctx context.Context, activity entity.Activity) (entity.Activity, error) {
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		return entity.Activity{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := r.Queries.WithTx(tx)

	var locationId *int32
	if activity.Location != nil {
		persistedLocationId, err := SaveLocation(ctx, qtx, *activity.Location)
		if err != nil {
			return entity.Activity{}, fmt.Errorf("save location: %w", err)
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
		return entity.Activity{}, fmt.Errorf("insert activity: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return entity.Activity{}, fmt.Errorf("commit tx: %w", err)
	}

	return r.GetActivityByID(ctx, activity.TripID, activityId)
}

func (r *ActivitiesRepo) UpdateActivity(ctx context.Context, activity entity.Activity) error {
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := r.Queries.WithTx(tx)

	existingLocationId, err := GetLocationIDOrNilByActivityID(ctx, qtx, activity.ID)
	if err != nil {
		return fmt.Errorf("get existing location id for activity [id=%d]: %w", activity.ID, err)
	}

	locationId, err := UpsertOrDeleteLocation(ctx, qtx, existingLocationId, activity.Location)
	if err != nil {
		return fmt.Errorf("upsert location [existing=%d]: %w", existingLocationId, err)
	}

	err = qtx.UpdateActivity(ctx, sqlc.UpdateActivityParams{
		ID:          activity.ID,
		LocationID:  locationId,
		Name:        activity.Name,
		Date:        activity.Date,
		Time:        activity.Time,
		Address:     activity.Address,
		Description: activity.Description,
		Price:       activity.Price,
	})
	if err != nil {
		return fmt.Errorf("update activity [id=%d]: %w", activity.ID, err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (r *ActivitiesRepo) DeleteActivity(ctx context.Context, tripID int32, activityID int32) error {
	return r.Queries.DeleteActivityByID(ctx, sqlc.DeleteActivityByIDParams{TripID: tripID, ID: activityID})
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
