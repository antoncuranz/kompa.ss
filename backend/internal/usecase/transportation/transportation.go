package transportation

import (
	"context"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
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
	})
	if err != nil {
		return entity.Transportation{}, err
	}

	return transportation, uc.saveGeoJson(ctx, transportation)
}

func (uc *UseCase) saveGeoJson(ctx context.Context, transportation entity.Transportation) error {
	featureCollection, err := uc.ors.LookupDirections(ctx, transportation.Origin, transportation.Destination)
	if err != nil {
		return fmt.Errorf("lookup directions: %w", err)
	}

	featureCollection.ExtraMembers = map[string]interface{}{"transportationType": transportation.Type}
	featureCollection.Append(featureWithProperties(transportation.Origin, transportation))
	featureCollection.Append(featureWithProperties(transportation.Destination, transportation))

	err = uc.repo.SaveGeoJson(ctx, transportation.ID, featureCollection)
	if err != nil {
		return fmt.Errorf("save geojson: %w", err)
	}
	return nil
}

func featureWithProperties(location entity.Location, transportation entity.Transportation) *geojson.Feature {
	feature := geojson.NewFeature(locationToPoint(location))

	feature.Properties["type"] = transportation.Type
	feature.Properties["name"] = "[TODO] Transportation.Name"
	feature.Properties["departureDateTime"] = transportation.DepartureDateTime
	feature.Properties["arrivalDateTime"] = transportation.ArrivalDateTime

	return feature
}

func locationToPoint(location entity.Location) orb.Point {
	return orb.Point{
		float64(location.Longitude),
		float64(location.Latitude),
	}
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

func (uc *UseCase) GetAllGeoJson(ctx context.Context, tripID int32) ([]geojson.FeatureCollection, error) {
	return uc.repo.GetAllGeoJson(ctx, tripID)
}
