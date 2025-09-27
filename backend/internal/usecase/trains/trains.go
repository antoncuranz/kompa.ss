package trains

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
	repo    repo.TransportationRepo
	dbVendo repo.DbVendoWebAPI
}

func New(r repo.TransportationRepo, a repo.DbVendoWebAPI) *UseCase {
	return &UseCase{
		repo:    r,
		dbVendo: a,
	}
}

func (uc *UseCase) RetrieveLocation(ctx context.Context, query string) (entity.TrainStation, error) {
	return uc.dbVendo.RetrieveLocation(ctx, query)
}

func (uc *UseCase) CreateTrainJourney(ctx context.Context, userID int32, tripID int32, journeyRequest request.TrainJourney) (entity.Transportation, error) {
	trainDetail, err := uc.dbVendo.RetrieveJourney(ctx, journeyRequest)
	if err != nil {
		return entity.Transportation{}, fmt.Errorf("failed to retrieve journey: %w", err)
	}

	firstLeg := trainDetail.Legs[0]
	lastLeg := trainDetail.Legs[len(trainDetail.Legs)-1]

	transportation, err := uc.repo.SaveTransportation(ctx, userID, entity.Transportation{
		TripID:            tripID,
		Type:              entity.TRAIN,
		Origin:            firstLeg.Origin.Location,
		Destination:       lastLeg.Destination.Location,
		DepartureDateTime: firstLeg.DepartureDateTime,
		ArrivalDateTime:   lastLeg.ArrivalDateTime,
		Price:             nil,
		TrainDetail:       &trainDetail,
	})
	if err != nil {
		return entity.Transportation{}, err
	}

	_, err = uc.retrieveAndPersistPolyline(ctx, userID, transportation.ID, trainDetail.RefreshToken)
	if err != nil {
		return entity.Transportation{}, fmt.Errorf("retrieve and process polyline: %w", err)
	}

	return transportation, nil
}

func (uc *UseCase) retrieveAndPersistPolyline(ctx context.Context, userID int32, transportationID int32, refreshToken string) (*geojson.FeatureCollection, error) {
	polylines, err := uc.dbVendo.RetrievePolylines(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("retrieve polyline: %w", err)
	}
	if len(polylines) == 0 {
		return nil, nil
	}

	featureCollection := geojson.NewFeatureCollection()
	featureCollection.ExtraMembers = map[string]interface{}{"transportationType": "TRAIN"}

	// start point
	featureCollection.Append(geojson.NewFeature(
		polylines[0].Features[0].Geometry.(orb.Point),
	))

	for _, polyline := range polylines {
		// line
		lineString := orb.LineString{}
		for _, feature := range polyline.Features {
			lineString = append(lineString, feature.Geometry.(orb.Point))
		}
		featureCollection.Append(geojson.NewFeature(lineString))

		// intermediate/end point
		featureCollection.Append(geojson.NewFeature(
			polyline.Features[len(polyline.Features)-1].Geometry.(orb.Point),
		))
	}

	err = uc.repo.SaveGeoJson(ctx, userID, transportationID, featureCollection)
	if err != nil {
		return nil, fmt.Errorf("save geojson: %w", err)
	}

	return featureCollection, nil
}
