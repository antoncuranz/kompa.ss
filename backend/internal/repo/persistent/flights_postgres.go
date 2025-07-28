package persistent

import (
	"backplate/internal/entity"
	"backplate/pkg/postgres"
	"backplate/pkg/sqlc"
	"context"
	"fmt"
	"time"
)

type FlightsRepo struct {
	*sqlc.Queries
}

func NewFlightsRepo(pg *postgres.Postgres) *FlightsRepo {
	return &FlightsRepo{sqlc.New(pg.Pool)}
}

func (r *FlightsRepo) GetFlights(ctx context.Context) ([]entity.Flight, error) {
	flights, err := r.Queries.GetFlights(ctx)
	if err != nil {
		return nil, err
	}

	var result []entity.Flight
	for _, flight := range flights {
		legs, err := r.Queries.GetFlightLegsByFlightID(ctx, flight.ID)
		if err != nil {
			return []entity.Flight{}, err
		}

		pnrs, err := r.Queries.GetPnrsByFlightID(ctx, flight.ID)
		if err != nil {
			return []entity.Flight{}, err
		}
		result = append(result, mapFlight(flight, legs, pnrs))
	}
	return result, nil
}

func (r *FlightsRepo) GetFlightByID(ctx context.Context, id int32) (entity.Flight, error) {
	flight, err := r.Queries.GetFlightByID(ctx, id)
	if err != nil {
		return entity.Flight{}, err
	}

	legs, err := r.Queries.GetFlightLegsByFlightID(ctx, flight.ID)
	if err != nil {
		return entity.Flight{}, err
	}

	pnrs, err := r.Queries.GetPnrsByFlightID(ctx, flight.ID)
	if err != nil {
		return entity.Flight{}, err
	}

	return mapFlight(flight, legs, pnrs), nil
}

// TODO: Move mapping to separate file

func mapFlight(flight sqlc.Flight, legs []sqlc.GetFlightLegsByFlightIDRow, pnrs []sqlc.Pnr) entity.Flight {
	return entity.Flight{
		ID:     flight.ID,
		TripID: flight.TripID,
		Price:  flight.Price,
		Legs:   mapFlightLegs(legs),
		PNRs:   mapPnrs(pnrs),
	}
}

func mapFlightLegs(legs []sqlc.GetFlightLegsByFlightIDRow) []entity.FlightLeg {
	var result []entity.FlightLeg
	for _, leg := range legs {
		result = append(result, mapFlightLeg(leg))
	}
	return result
}

func mapFlightLeg(leg sqlc.GetFlightLegsByFlightIDRow) entity.FlightLeg {
	format := "2006-01-02 15:04-07:00"
	departureTime, err := time.Parse(format, leg.FlightLeg.DepartureTime)
	if err != nil {
		fmt.Printf("Error parsing timestamp: %s\n", leg.FlightLeg.DepartureTime)
	}
	arrivalTime, err := time.Parse(format, leg.FlightLeg.ArrivalTime)
	if err != nil {
		fmt.Printf("Error parsing timestamp: %s\n", leg.FlightLeg.ArrivalTime)
	}

	return entity.FlightLeg{
		ID:            leg.FlightLeg.ID,
		Origin:        mapAirport(leg.Airport),
		Destination:   mapAirport(leg.Airport_2),
		Airline:       leg.FlightLeg.Airline,
		FlightNumber:  leg.FlightLeg.FlightNumber,
		DepartureTime: departureTime,
		ArrivalTime:   arrivalTime,
		Aircraft:      leg.FlightLeg.Aircraft,
	}
}

func mapAirport(airport sqlc.Airport) entity.Airport {
	return entity.Airport{
		Iata:         airport.Iata,
		Name:         airport.Name,
		Municipality: airport.Municipality,
		Location:     airport.Location,
	}
}

func mapPnrs(pnrs []sqlc.Pnr) []entity.PNR {
	var result []entity.PNR
	for _, pnr := range pnrs {
		result = append(result, mapPnr(pnr))
	}
	return result
}

func mapPnr(pnr sqlc.Pnr) entity.PNR {
	return entity.PNR{
		ID:      pnr.ID,
		Airline: pnr.Airline,
		PNR:     pnr.Pnr,
	}
}
