package transportation

import (
	"context"
	"fmt"
	"github.com/paulmach/orb/geojson"
	"kompass/internal/entity"
	"kompass/internal/repo"
)

type UseCase struct {
	repo repo.TransportationRepo
}

func New(r repo.TransportationRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) GetAllTransportation(ctx context.Context, userID int32, tripID int32) ([]entity.Transportation, error) {
	allTransportation, err := uc.repo.GetAllTransportation(ctx, userID, tripID)
	if err != nil {
		return nil, fmt.Errorf("get all transportation: %w", err)
	}

	return allTransportation, nil
}

func (uc *UseCase) GetTransportationByID(ctx context.Context, userID int32, tripID int32, transportationID int32) (entity.Transportation, error) {
	transportation, err := uc.repo.GetTransportationByID(ctx, userID, tripID, transportationID)
	if err != nil {
		return entity.Transportation{}, fmt.Errorf("get transportation [id=%d]: %w", transportationID, err)
	}

	return transportation, nil
}

func (uc *UseCase) DeleteTransportation(ctx context.Context, userID int32, tripID int32, transportationID int32) error {
	return uc.repo.DeleteTransportation(ctx, userID, tripID, transportationID)
}

func (uc *UseCase) GetAllGeoJson(ctx context.Context, userID int32, tripID int32) ([]geojson.FeatureCollection, error) {
	return uc.repo.GetAllGeoJson(ctx, userID, tripID)
}
