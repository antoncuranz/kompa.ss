package persistent

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"travel-planner/internal/entity"
	"travel-planner/pkg/postgres"
	"travel-planner/pkg/sqlc"
)

type FlightsRepo struct {
	Db      *pgxpool.Pool
	Queries *sqlc.Queries
}

func NewFlightsRepo(pg *postgres.Postgres) *FlightsRepo {
	return &FlightsRepo{
		pg.Pool,
		sqlc.New(pg.Pool),
	}
}

func (r *FlightsRepo) GetFlights(ctx context.Context) ([]entity.Flight, error) {
	flights, err := r.Queries.GetFlights(ctx)
	if err != nil {
		return nil, err
	}

	result := []entity.Flight{}
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

func (r *FlightsRepo) SaveFlight(ctx context.Context, flight entity.Flight) (entity.Flight, error) {

	tx, err := r.Db.Begin(ctx)
	if err != nil {
		return entity.Flight{}, err
	}
	defer tx.Rollback(ctx)
	qtx := r.Queries.WithTx(tx)

	flightId, err := qtx.InsertFlight(ctx, sqlc.InsertFlightParams{
		TripID: flight.TripID,
		Price:  flight.Price,
	})
	if err != nil {
		return entity.Flight{}, err
	}

	airportSet := make(map[string]entity.Airport)
	for _, leg := range flight.Legs {
		airportSet[leg.Origin.Iata] = leg.Origin
		airportSet[leg.Destination.Iata] = leg.Destination
	}

	for iata := range airportSet {
		airport := airportSet[iata]
		err := qtx.InsertAirport(ctx, sqlc.InsertAirportParams{
			Iata:         airport.Iata,
			Name:         airport.Name,
			Municipality: airport.Municipality,
			Location:     airport.Location,
		})
		if err != nil {
			return entity.Flight{}, err
		}
	}

	for _, leg := range flight.Legs {
		_, err := qtx.InsertFlightLeg(ctx, sqlc.InsertFlightLegParams{
			FlightID:      flightId,
			Origin:        leg.Origin.Iata,
			Destination:   leg.Destination.Iata,
			Airline:       leg.Airline,
			FlightNumber:  leg.FlightNumber,
			DepartureTime: leg.DepartureTime,
			ArrivalTime:   leg.ArrivalTime,
			Aircraft:      leg.Aircraft,
		})
		if err != nil {
			return entity.Flight{}, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return entity.Flight{}, err
	}

	return r.GetFlightByID(ctx, flightId)
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
	result := []entity.FlightLeg{}
	for _, leg := range legs {
		result = append(result, mapFlightLeg(leg))
	}
	return result
}

func mapFlightLeg(leg sqlc.GetFlightLegsByFlightIDRow) entity.FlightLeg {
	//format := "2006-01-02 15:04-07:00"
	//departureTime, err := time.Parse(format, leg.FlightLeg.DepartureTime)
	//if err != nil {
	//	fmt.Printf("Error parsing timestamp: %s\n", leg.FlightLeg.DepartureTime)
	//}
	//arrivalTime, err := time.Parse(format, leg.FlightLeg.ArrivalTime)
	//if err != nil {
	//	fmt.Printf("Error parsing timestamp: %s\n", leg.FlightLeg.ArrivalTime)
	//}

	return entity.FlightLeg{
		ID:            leg.FlightLeg.ID,
		Origin:        mapAirport(leg.Airport),
		Destination:   mapAirport(leg.Airport_2),
		Airline:       leg.FlightLeg.Airline,
		FlightNumber:  leg.FlightLeg.FlightNumber,
		DepartureTime: leg.FlightLeg.DepartureTime,
		ArrivalTime:   leg.FlightLeg.ArrivalTime,
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
	result := []entity.PNR{}
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
