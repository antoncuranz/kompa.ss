package trains

import (
	"context"
	"fmt"
	"kompass/internal/controller/http/v1/request"
	"kompass/internal/entity"
	"kompass/internal/repo"
)

type UseCase struct {
	dbVendo repo.DbVendoWebAPI
}

func New(api repo.DbVendoWebAPI) *UseCase {
	return &UseCase{
		dbVendo: api,
	}
}

func (uc *UseCase) LookupTrainStation(ctx context.Context, query string) (entity.TrainStation, error) {
	return uc.dbVendo.LookupTrainStation(ctx, query)
}

func (uc *UseCase) FindTrainJourney(ctx context.Context, journeyRequest request.TrainJourney) (entity.TrainDetail, error) {
	trainDetail, err := uc.dbVendo.RetrieveJourney(ctx, journeyRequest)
	if err != nil {
		return entity.TrainDetail{}, fmt.Errorf("failed to retrieve journey: %w", err)
	}

	return trainDetail, nil
	//_, err = uc.createGeoJson(ctx, transportation)
}
