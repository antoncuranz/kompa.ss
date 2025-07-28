package flights

import (
	"backplate/internal/entity"
	"backplate/internal/repo"
	"context"
)

type UseCase struct {
	repo repo.FlightsRepo
}

func New(r repo.FlightsRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) GetFlightByID(ctx context.Context, id int32) (entity.Flight, error) {
	return uc.repo.GetFlightByID(ctx, id)
}

func (uc *UseCase) GetFlights(ctx context.Context) ([]entity.Flight, error) {
	return uc.repo.GetFlights(ctx)
}
