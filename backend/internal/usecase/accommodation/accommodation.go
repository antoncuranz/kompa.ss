package accommodation

import (
	"context"
	"travel-planner/internal/entity"
	"travel-planner/internal/repo"
)

type UseCase struct {
	repo repo.AccommodationRepo
}

func New(r repo.AccommodationRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) GetAccommodationByID(ctx context.Context, id int32) (entity.Accommodation, error) {
	return uc.repo.GetAccommodationByID(ctx, id)
}

func (uc *UseCase) GetAllAccommodation(ctx context.Context) ([]entity.Accommodation, error) {
	return uc.repo.GetAllAccommodation(ctx)
}
