package activities

import (
	"context"
	"travel-planner/internal/controller/http/v1/request"
	"travel-planner/internal/entity"
	"travel-planner/internal/repo"
)

type UseCase struct {
	repo repo.ActivitiesRepo
}

func New(r repo.ActivitiesRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) GetActivityByID(ctx context.Context, id int32) (entity.Activity, error) {
	return uc.repo.GetActivityByID(ctx, id)
}

func (uc *UseCase) GetActivities(ctx context.Context) ([]entity.Activity, error) {
	return uc.repo.GetActivities(ctx)
}

func (uc *UseCase) CreateActivity(ctx context.Context, activity request.Activity) (entity.Activity, error) {
	return uc.repo.SaveActivity(ctx, entity.Activity{
		TripID:      activity.TripID,
		Name:        activity.Name,
		Date:        activity.Date,
		Time:        activity.Time,
		Description: activity.Description,
		Address:     activity.Address,
		Location:    activity.Location,
		Price:       activity.Price,
	})
}
