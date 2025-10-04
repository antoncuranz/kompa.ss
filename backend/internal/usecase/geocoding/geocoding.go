package geocoding

import (
	"context"
	"fmt"
	"github.com/valyala/fasthttp"
	"kompass/internal/entity"
	"kompass/internal/repo"
	"kompass/internal/usecase"
)

type UseCase struct {
	trains usecase.Trains
	orm    repo.OpenRouteServiceWebAPI
}

func New(trains usecase.Trains, orm repo.OpenRouteServiceWebAPI) *UseCase {
	return &UseCase{
		trains: trains,
		orm:    orm,
	}
}

func (uc *UseCase) LookupLocation(ctx *fasthttp.RequestCtx, query string) (entity.GeocodeLocation, error) {
	location, err := uc.orm.LookupLocation(ctx, query)
	if err != nil {
		return entity.GeocodeLocation{}, fmt.Errorf("lookup location: %w", err)
	}

	return location, nil
}

func (uc *UseCase) LookupTrainStation(ctx context.Context, query string) (entity.TrainStation, error) {
	station, err := uc.trains.LookupTrainStation(ctx, query)
	if err != nil {
		return entity.TrainStation{}, fmt.Errorf("lookup trainstation: %w", err)
	}

	return station, nil
}
