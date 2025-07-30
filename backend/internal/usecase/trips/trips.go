package trips

import (
	"context"
	"travel-planner/internal/entity"
	"travel-planner/internal/repo"
)

type UseCase struct {
	repo repo.TripsRepo
}

func New(r repo.TripsRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) GetTripByID(ctx context.Context, id int32) (entity.Trip, error) {
	return uc.repo.GetTripByID(ctx, id)
}

func (uc *UseCase) GetTrips(ctx context.Context) ([]entity.Trip, error) {
	return uc.repo.GetTrips(ctx)
}
