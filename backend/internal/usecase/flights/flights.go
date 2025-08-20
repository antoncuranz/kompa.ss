package flights

import (
	"context"
	"fmt"
	"kompass/internal/controller/http/v1/request"
	"kompass/internal/entity"
	"kompass/internal/repo"
)

type UseCase struct {
	repo        repo.FlightsRepo
	aerodatabox repo.AerodataboxWebAPI
}

func New(r repo.FlightsRepo, a repo.AerodataboxWebAPI) *UseCase {
	return &UseCase{
		repo:        r,
		aerodatabox: a,
	}
}

func (uc *UseCase) GetFlights(ctx context.Context, tripID int32) ([]entity.Flight, error) {
	return uc.repo.GetFlights(ctx, tripID)
}

func (uc *UseCase) GetFlightByID(ctx context.Context, tripID int32, flightID int32) (entity.Flight, error) {
	return uc.repo.GetFlightByID(ctx, tripID, flightID)
}

func (uc *UseCase) CreateFlight(ctx context.Context, tripID int32, flight request.Flight) (entity.Flight, error) {
	flightLegs, err := uc.retrieveFlightLegs(ctx, flight)
	if err != nil {
		return entity.Flight{}, err
	}

	return uc.repo.SaveFlight(ctx, entity.Flight{
		TripID: tripID,
		Legs:   flightLegs,
		PNRs:   flight.PNRs,
		Price:  flight.Price,
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

func (uc *UseCase) DeleteFlight(ctx context.Context, tripID int32, flightID int32) error {
	return uc.repo.DeleteFlight(ctx, tripID, flightID)
}
