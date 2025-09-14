package flights

import (
	"context"
	"fmt"
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

func (uc *UseCase) CreateFlight(ctx context.Context, tripID int32, flight request.Flight) (entity.Transportation, error) {
	flightLegs, err := uc.retrieveFlightLegs(ctx, flight)
	if err != nil {
		return entity.Transportation{}, err
	}

	firstLeg := flightLegs[0]
	lastLeg := flightLegs[len(flightLegs)-1]

	return uc.repo.SaveTransportation(ctx, entity.Transportation{
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

func (uc *UseCase) UpdateFlight(ctx context.Context, tripID int32, flightID int32, flight request.Flight) error {
	return fmt.Errorf("Not yet implemented")
}
