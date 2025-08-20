package activities

import (
	"context"
	"kompass/internal/controller/http/v1/request"
	"kompass/internal/entity"
	"kompass/internal/repo"
)

type UseCase struct {
	repo repo.ActivitiesRepo
}

func New(r repo.ActivitiesRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) GetActivities(ctx context.Context, tripID int32) ([]entity.Activity, error) {
	return uc.repo.GetActivities(ctx, tripID)
}

func (uc *UseCase) GetActivityByID(ctx context.Context, tripID int32, id int32) (entity.Activity, error) {
	return uc.repo.GetActivityByID(ctx, tripID, id)
}

func (uc *UseCase) CreateActivity(ctx context.Context, tripID int32, activity request.Activity) (entity.Activity, error) {
	return uc.repo.SaveActivity(ctx, entity.Activity{
		TripID:      tripID,
		Name:        activity.Name,
		Date:        activity.Date,
		Time:        activity.Time,
		Description: activity.Description,
		Address:     activity.Address,
		Location:    activity.Location,
		Price:       activity.Price,
	})
}

func (uc *UseCase) UpdateActivity(ctx context.Context, tripID int32, activityID int32, activity request.Activity) error {
	return uc.repo.UpdateActivity(ctx, entity.Activity{
		ID:          activityID,
		TripID:      tripID,
		Name:        activity.Name,
		Date:        activity.Date,
		Time:        activity.Time,
		Description: activity.Description,
		Address:     activity.Address,
		Location:    activity.Location,
		Price:       activity.Price,
	})
}

func (uc *UseCase) DeleteActivity(ctx context.Context, tripID int32, activityId int32) error {
	return uc.repo.DeleteActivity(ctx, tripID, activityId)
}
