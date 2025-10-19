package flights

import (
	"context"
	"fmt"
	"kompass/internal/controller/http/v1/request"
	"kompass/internal/entity"
	"kompass/internal/repo"
	"sort"
	"strings"

	"cloud.google.com/go/civil"
)

type UseCase struct {
	transportationRepo repo.TransportationRepo
	flightsRepo        repo.FlightsRepo
	flightsApi         repo.FlightInformationWebAPI
}

func New(transportationRepo repo.TransportationRepo, flightsRepo repo.FlightsRepo, a repo.FlightInformationWebAPI) *UseCase {
	return &UseCase{
		transportationRepo: transportationRepo,
		flightsRepo:        flightsRepo,
		flightsApi:         a,
	}
}

func (uc *UseCase) CreateFlight(ctx context.Context, tripID int32, flight request.Flight) (entity.Transportation, error) {
	flightLegs, err := uc.retrieveFlightLegs(ctx, flight)
	if err != nil {
		return entity.Transportation{}, err
	}

	sortByDepartureDate(flightLegs)
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

func (uc *UseCase) UpdateFlight(ctx context.Context, tripID int32, flightID int32) error {
	// FIXME: use tripID in query!
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

	sortByDepartureDate(flightLegs)
	firstLeg := flightLegs[0]
	lastLeg := flightLegs[len(flightLegs)-1]
	transportation.DepartureDateTime = firstLeg.DepartureDateTime
	transportation.ArrivalDateTime = lastLeg.ArrivalDateTime
	transportation.Origin = firstLeg.Origin.Location
	transportation.Destination = lastLeg.Destination.Location
	// TODO: update PNRs and price

	updated, err := uc.transportationRepo.SaveTransportation(ctx, transportation)
	if err != nil {
		return fmt.Errorf("update transportation: %w", err)
	}

	return uc.saveGeoJson(ctx, updated)
}

func (uc *UseCase) retrieveFlightLegs(ctx context.Context, flight request.Flight) ([]entity.FlightLeg, error) {
	legs := []entity.FlightLeg{}
	for _, leg := range flight.Legs {
		flightLeg, err := uc.flightsApi.RetrieveFlightLeg(ctx, leg.Date, leg.FlightNumber, leg.OriginAirport)
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
		flightLeg, err := uc.flightsApi.RetrieveFlightLeg(ctx, getFlightDate(leg), flightNumber, &leg.Origin.Iata)
		if err != nil {
			return []entity.FlightLeg{}, err
		}
		flightLeg.ID = leg.ID
		legs = append(legs, flightLeg)
	}

	return legs, nil
}

func sortByDepartureDate(legs []entity.FlightLeg) {
	sort.Slice(legs, func(i, j int) bool {
		return legs[i].DepartureDateTime.Compare(legs[j].DepartureDateTime) < 0
	})
}

func getFlightDate(leg entity.FlightLeg) civil.Date {
	if leg.AmadeusFlightDate != nil {
		return *leg.AmadeusFlightDate
	}
	return leg.DepartureDateTime.Date
}
