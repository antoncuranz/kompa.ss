package transportation

import (
	"context"
	"fmt"
	"kompass/internal/controller/http/v1/request"
	"kompass/internal/entity"
	"kompass/internal/repo"
)

type UseCase struct {
	repo repo.TransportationRepo
	ors  repo.OpenRouteServiceWebAPI
}

func New(r repo.TransportationRepo, ors repo.OpenRouteServiceWebAPI) *UseCase {
	return &UseCase{
		repo: r,
		ors:  ors,
	}
}

func (uc *UseCase) CreateTransportation(ctx context.Context, tripID int32, request request.Transportation) (entity.Transportation, error) {
	transportation, err := uc.repo.SaveTransportation(ctx, entity.Transportation{
		TripID:            tripID,
		Type:              request.Type,
		Origin:            request.Origin,
		Destination:       request.Destination,
		DepartureDateTime: request.DepartureDateTime,
		ArrivalDateTime:   request.ArrivalDateTime,
		Price:             request.Price,
		GenericDetail: &entity.GenericDetail{
			Name:               request.Name,
			OriginAddress:      request.OriginAddress,
			DestinationAddress: request.DestinationAddress,
		},
	})
	if err != nil {
		return entity.Transportation{}, err
	}

	return transportation, uc.saveGeoJson(ctx, transportation)
}

func (uc *UseCase) UpdateTransportation(ctx context.Context, tripID int32, transportationID int32, request request.Transportation) (entity.Transportation, error) {
	transportation, err := uc.repo.SaveTransportation(ctx, entity.Transportation{
		ID:                transportationID,
		TripID:            tripID,
		Type:              request.Type,
		Origin:            request.Origin,
		Destination:       request.Destination,
		DepartureDateTime: request.DepartureDateTime,
		ArrivalDateTime:   request.ArrivalDateTime,
		Price:             request.Price,
		GenericDetail: &entity.GenericDetail{
			Name:               request.Name,
			OriginAddress:      request.OriginAddress,
			DestinationAddress: request.DestinationAddress,
		},
	})
	if err != nil {
		return entity.Transportation{}, err
	}

	return transportation, uc.saveGeoJson(ctx, transportation)
}

func (uc *UseCase) GetAllTransportation(ctx context.Context, tripID int32) ([]entity.Transportation, error) {
	allTransportation, err := uc.repo.GetAllTransportation(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("get all transportation: %w", err)
	}

	return allTransportation, nil
}

func (uc *UseCase) GetTransportationByID(ctx context.Context, tripID int32, transportationID int32) (entity.Transportation, error) {
	transportation, err := uc.repo.GetTransportationByID(ctx, tripID, transportationID)
	if err != nil {
		return entity.Transportation{}, fmt.Errorf("get transportation [id=%d]: %w", transportationID, err)
	}

	return transportation, nil
}

func (uc *UseCase) DeleteTransportation(ctx context.Context, tripID int32, transportationID int32) error {
	return uc.repo.DeleteTransportation(ctx, tripID, transportationID)
}
