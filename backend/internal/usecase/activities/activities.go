package activities

import (
	"backplate/internal/entity"
	"backplate/internal/repo"
	"context"
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
