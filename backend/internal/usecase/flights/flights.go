package flights

import (
	"context"
	"travel-planner/internal/controller/http/v1/request"
	"travel-planner/internal/entity"
	"travel-planner/internal/repo"
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

func (uc *UseCase) GetFlightByID(ctx context.Context, id int32) (entity.Flight, error) {
	return uc.repo.GetFlightByID(ctx, id)
}

func (uc *UseCase) GetFlights(ctx context.Context) ([]entity.Flight, error) {
	return uc.repo.GetFlights(ctx)
}

func (uc *UseCase) CreateFlight(ctx context.Context, flight request.Flight) (entity.Flight, error) {
	flightLegs, err := uc.retrieveFlightLegs(ctx, flight)
	if err != nil {
		return entity.Flight{}, err
	}

	return uc.repo.SaveFlight(ctx, entity.Flight{
		TripID: 1,
		Legs:   flightLegs,
		PNRs:   flight.PNRs,
		Price:  nil,
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
