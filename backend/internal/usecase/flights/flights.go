package flights

import (
	"context"
	"kompass/internal/controller/http/v1/request"
	"kompass/internal/entity"
	"kompass/internal/repo"
	"sort"
	"strings"

	"cloud.google.com/go/civil"
)

type UseCase struct {
	flightsApi repo.FlightInformationWebAPI
}

func New(a repo.FlightInformationWebAPI) *UseCase {
	return &UseCase{
		flightsApi: a,
	}
}

func (uc *UseCase) FindFlight(ctx context.Context, flight request.Flight) (entity.FlightDetail, error) {
	flightLegs, err := uc.retrieveFlightLegs(ctx, flight)
	if err != nil {
		return entity.FlightDetail{}, err
	}

	sortByDepartureDate(flightLegs)

	flightDetail := entity.FlightDetail{
		Legs: flightLegs,
		PNRs: flight.PNRs,
	}

	//uc.createGeoJson(flightDetail)

	return flightDetail, nil
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
