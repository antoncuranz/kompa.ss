package flights

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
	repo        repo.TransportationRepo
	aerodatabox repo.AerodataboxWebAPI
}

func New(r repo.TransportationRepo, a repo.AerodataboxWebAPI) *UseCase {
	return &UseCase{
		repo:        r,
		aerodatabox: a,
	}
}

func (uc *UseCase) CreateFlight(ctx context.Context, userID int32, tripID int32, flight request.Flight) (entity.Transportation, error) {
	flightLegs, err := uc.retrieveFlightLegs(ctx, flight)
	if err != nil {
		return entity.Transportation{}, err
	}

	firstLeg := flightLegs[0]
	lastLeg := flightLegs[len(flightLegs)-1]

	transportation, err := uc.repo.SaveTransportation(ctx, userID, entity.Transportation{
		TripID:            tripID,
		Type:              entity.PLANE,
		Origin:            firstLeg.Origin.Location,
		Destination:       lastLeg.Destination.Location,
		DepartureDateTime: firstLeg.DepartureDateTime,
		ArrivalDateTime:   lastLeg.ArrivalDateTime,
		Price:             flight.Price,
		FlightDetail: &entity.FlightDetail{
			Legs: flightLegs,
			PNRs: flight.PNRs,
		},
	})
	if err != nil {
		return entity.Transportation{}, err
	}

	return transportation, uc.SaveGeoJson(ctx, userID, transportation)
}

func (uc *UseCase) SaveGeoJson(ctx context.Context, userID int32, transportation entity.Transportation) error {

	legs := transportation.FlightDetail.Legs

	featureCollection := geojson.NewFeatureCollection()
	featureCollection.ExtraMembers = map[string]interface{}{"transportationType": "PLANE"}

	// start point
	featureCollection.Append(geojson.NewFeature(
		locationToPoint(legs[0].Origin.Location),
	))

	for _, leg := range legs {
		// line
		featureCollection.Append(geojson.NewFeature(
			orb.LineString{
				locationToPoint(leg.Origin.Location),
				locationToPoint(leg.Destination.Location),
			},
		))

		// intermediate/end point
		featureCollection.Append(geojson.NewFeature(
			locationToPoint(leg.Destination.Location),
		))
	}

	err := uc.repo.SaveGeoJson(ctx, userID, transportation.ID, featureCollection)
	if err != nil {
		return fmt.Errorf("save geojson: %w", err)
	}
	return nil
}

func locationToPoint(location entity.Location) orb.Point {
	return orb.Point{
		float64(location.Longitude),
		float64(location.Latitude),
	}
}

func (uc *UseCase) retrieveFlightLegs(ctx context.Context, flight request.Flight) ([]entity.FlightLeg, error) {
	legs := []entity.FlightLeg{}
	for _, leg := range flight.Legs {
		flightLeg, err := uc.aerodatabox.RetrieveFlightLeg(ctx, leg.Date, leg.FlightNumber, leg.OriginAirport)
		if err != nil {
			return []entity.FlightLeg{}, err
		}
		legs = append(legs, flightLeg)
	}

	return legs, nil
}

func (uc *UseCase) UpdateFlight(ctx context.Context, userID int32, tripID int32, flightID int32, flight request.Flight) error {
	return fmt.Errorf("Not yet implemented")
}
