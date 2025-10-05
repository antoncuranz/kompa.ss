package trains

import (
	"context"
	"fmt"
	"kompass/internal/controller/http/v1/request"
	"kompass/internal/entity"
	"kompass/internal/repo"
)

type UseCase struct {
	repo    repo.TransportationRepo
	dbVendo repo.DbVendoWebAPI
}

func New(r repo.TransportationRepo, a repo.DbVendoWebAPI) *UseCase {
	return &UseCase{
		repo:    r,
		dbVendo: a,
	}
}

func (uc *UseCase) LookupTrainStation(ctx context.Context, query string) (entity.TrainStation, error) {
	return uc.dbVendo.LookupTrainStation(ctx, query)
}

func (uc *UseCase) CreateTrainJourney(ctx context.Context, tripID int32, journeyRequest request.TrainJourney) (entity.Transportation, error) {
	trainDetail, err := uc.dbVendo.RetrieveJourney(ctx, journeyRequest)
	if err != nil {
		return entity.Transportation{}, fmt.Errorf("failed to retrieve journey: %w", err)
	}

	firstLeg := trainDetail.Legs[0]
	lastLeg := trainDetail.Legs[len(trainDetail.Legs)-1]

	transportation, err := uc.repo.SaveTransportation(ctx, entity.Transportation{
		TripID:            tripID,
		Type:              entity.TRAIN,
		Origin:            firstLeg.Origin.Location,
		Destination:       lastLeg.Destination.Location,
		DepartureDateTime: firstLeg.DepartureDateTime,
		ArrivalDateTime:   lastLeg.ArrivalDateTime,
		Price:             journeyRequest.Price,
		TrainDetail:       &trainDetail,
	})
	if err != nil {
		return entity.Transportation{}, err
	}

	_, err = uc.saveGeoJson(ctx, transportation)
	if err != nil {
		return entity.Transportation{}, fmt.Errorf("retrieve and process polyline: %w", err)
	}

	return transportation, nil
}
