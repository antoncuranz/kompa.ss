package flights

import (
	"context"
	"fmt"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"kompass/internal/controller/http/v1/request"
	"kompass/internal/entity"
	"kompass/internal/repo"
	"strings"
)

type UseCase struct {
	transportationRepo repo.TransportationRepo
	flightsRepo        repo.FlightsRepo
	aerodatabox        repo.AerodataboxWebAPI
}

func New(transportationRepo repo.TransportationRepo, flightsRepo repo.FlightsRepo, a repo.AerodataboxWebAPI) *UseCase {
	return &UseCase{
		transportationRepo: transportationRepo,
		flightsRepo:        flightsRepo,
		aerodatabox:        a,
	}
}

func (uc *UseCase) CreateFlight(ctx context.Context, tripID int32, flight request.Flight) (entity.Transportation, error) {
	flightLegs, err := uc.retrieveFlightLegs(ctx, flight)
	if err != nil {
		return entity.Transportation{}, err
	}

	firstLeg := flightLegs[0]
	lastLeg := flightLegs[len(flightLegs)-1]

	transportation, err := uc.transportationRepo.SaveTransportation(ctx, entity.Transportation{
		TripID:            tripID,
		Type:              entity.FLIGHT,
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

	return transportation, uc.saveGeoJson(ctx, transportation)
}

func (uc *UseCase) saveGeoJson(ctx context.Context, transportation entity.Transportation) error {

	legs := transportation.FlightDetail.Legs

	featureCollection := geojson.NewFeatureCollection()
	featureCollection.ExtraMembers = map[string]interface{}{"transportationType": "FLIGHT"}

	airportByIata := map[string]entity.Airport{}
	legsByAirport := map[string][]entity.FlightLeg{}

	for _, leg := range legs {
		featureCollection.Append(geojson.NewFeature(
			orb.LineString{
				locationToPoint(leg.Origin.Location),
				locationToPoint(leg.Destination.Location),
			},
		))

		airportByIata[leg.Origin.Iata] = leg.Origin
		airportByIata[leg.Destination.Iata] = leg.Destination
		legsByAirport[leg.Origin.Iata] = append(legsByAirport[leg.Origin.Iata], leg)
		legsByAirport[leg.Destination.Iata] = append(legsByAirport[leg.Destination.Iata], leg)
	}

	from := legs[0].Origin.Municipality
	to := legs[len(legs)-1].Destination.Municipality

	for iata, legs := range legsByAirport {
		location := airportByIata[iata].Location
		featureCollection.Append(featureWithProperties(from, to, location, legs))
	}

	err := uc.transportationRepo.SaveGeoJson(ctx, transportation.ID, featureCollection)
	if err != nil {
		return fmt.Errorf("save geojson: %w", err)
	}
	return nil
}

func featureWithProperties(fromMunicipality string, toMunicipality string, location entity.Location, legs []entity.FlightLeg) *geojson.Feature {
	feature := geojson.NewFeature(locationToPoint(location))

	feature.Properties["type"] = "FLIGHT"
	feature.Properties["fromMunicipality"] = fromMunicipality
	feature.Properties["toMunicipality"] = toMunicipality

	var legProperties []map[string]interface{}
	for _, leg := range legs {
		legProperties = append(legProperties, map[string]interface{}{
			"flightNumber":      leg.FlightNumber,
			"departureDateTime": leg.DepartureDateTime,
			"arrivalDateTime":   leg.ArrivalDateTime,
			"fromIata":          leg.Origin.Iata,
			"toIata":            leg.Destination.Iata,
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

func (uc *UseCase) retrieveFlightLegsUpdate(ctx context.Context, flight entity.FlightDetail) ([]entity.FlightLeg, error) {
	legs := []entity.FlightLeg{}
	for _, leg := range flight.Legs {
		flightNumber := strings.ReplaceAll(leg.FlightNumber, " ", "")
		flightLeg, err := uc.aerodatabox.RetrieveFlightLeg(ctx, leg.DepartureDateTime.Date, flightNumber, &leg.Origin.Iata)
		if err != nil {
			return []entity.FlightLeg{}, err
		}
		flightLeg.ID = leg.ID
		legs = append(legs, flightLeg)
	}

	return legs, nil
}

func (uc *UseCase) UpdateFlight(ctx context.Context, tripID int32, flightID int32) error {
	flightDetail, err := uc.flightsRepo.GetFlightDetail(ctx, flightID)
	if err != nil {
		return fmt.Errorf("get flight detail [id=%d]: %w", flightID, err)
	}

	flightLegs, err := uc.retrieveFlightLegsUpdate(ctx, flightDetail)
	if err != nil {
		return fmt.Errorf("retrieve flight leg information: %w", err)
	}

	err = uc.flightsRepo.UpdateFlightLegs(ctx, flightLegs)
	if err != nil {
		return fmt.Errorf("update flight leg [id=%d]: %w", flightID, err)
	}

	transportation, err := uc.transportationRepo.GetTransportationByID(ctx, tripID, flightID)
	if err != nil {
		return err
	}

	return uc.saveGeoJson(ctx, transportation)
}
