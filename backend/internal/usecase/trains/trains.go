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

func (uc *UseCase) saveGeoJson(ctx context.Context, transportation entity.Transportation) (*geojson.FeatureCollection, error) {
	legs := transportation.TrainDetail.Legs

	polylines, err := uc.dbVendo.RetrievePolylines(ctx, transportation.TrainDetail.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("retrieve polyline: %w", err)
	}
	if len(polylines) == 0 {
		return nil, nil
	}

	featureCollection := geojson.NewFeatureCollection()
	featureCollection.ExtraMembers = map[string]interface{}{"transportationType": "TRAIN"}

	stationByID := map[string]entity.TrainStation{}
	legsByStation := map[string][]entity.TrainLeg{}

	for _, polyline := range polylines {
		lineString := orb.LineString{}
		for _, feature := range polyline.Features {
			lineString = append(lineString, feature.Geometry.(orb.Point))
		}
		featureCollection.Append(geojson.NewFeature(lineString))
	}

	for _, leg := range legs {
		stationByID[leg.Origin.ID] = leg.Origin
		stationByID[leg.Destination.ID] = leg.Destination
		legsByStation[leg.Origin.ID] = append(legsByStation[leg.Origin.ID], leg)
		legsByStation[leg.Destination.ID] = append(legsByStation[leg.Destination.ID], leg)
	}

	from := legs[0].Origin.Name
	to := legs[len(legs)-1].Destination.Name

	for stationID, legs := range legsByStation {
		location := stationByID[stationID].Location
		featureCollection.Append(featureWithProperties(from, to, location, legs))
	}

	err = uc.repo.SaveGeoJson(ctx, transportation.ID, featureCollection)
	if err != nil {
		return nil, fmt.Errorf("save geojson: %w", err)
	}

	return featureCollection, nil
}

func featureWithProperties(fromMunicipality string, toMunicipality string, location entity.Location, legs []entity.TrainLeg) *geojson.Feature {
	feature := geojson.NewFeature(locationToPoint(location))

	feature.Properties["type"] = "TRAIN"
	feature.Properties["fromMunicipality"] = fromMunicipality
	feature.Properties["toMunicipality"] = toMunicipality

	var legProperties []map[string]interface{}
	for _, leg := range legs {
		legProperties = append(legProperties, map[string]interface{}{
			"lineName":          leg.LineName,
			"departureDateTime": leg.DepartureDateTime,
			"arrivalDateTime":   leg.ArrivalDateTime,
			"fromStation":       leg.Origin.Name,
			"toStation":         leg.Destination.Name,
		})
	}
	feature.Properties["legs"] = legProperties

	return feature
}

func locationToPoint(location entity.Location) orb.Point {
	return orb.Point{
		float64(location.Longitude),
		float64(location.Latitude),
	}
}
