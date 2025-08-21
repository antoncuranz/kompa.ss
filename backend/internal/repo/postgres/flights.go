package postgres

import (
	"context"
	"kompass/internal/entity"
	"kompass/internal/repo/postgres/converter"
	"kompass/pkg/postgres"
	"kompass/pkg/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FlightsRepo struct {
	Db      *pgxpool.Pool
	Queries *sqlc.Queries
	c       converter.FlightConverter
}

func NewFlightsRepo(pg *postgres.Postgres) *FlightsRepo {
	return &FlightsRepo{
		pg.Pool,
		sqlc.New(pg.Pool),
		&converter.FlightConverterImpl{},
	}
}

func (r *FlightsRepo) GetFlights(ctx context.Context, tripID int32) ([]entity.Flight, error) {
	flights, err := r.Queries.GetFlights(ctx, tripID)
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
		result = append(result, r.c.ConvertFlight(converter.ConvertFlightParams{Flight: flight, Legs: legs, PNRs: pnrs}))
	}
	return result, nil
}

func (r *FlightsRepo) GetFlightByID(ctx context.Context, tripID int32, flightID int32) (entity.Flight, error) {
	flight, err := r.Queries.GetFlightByID(ctx, sqlc.GetFlightByIDParams{TripID: tripID, ID: flightID})
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

	return r.c.ConvertFlight(converter.ConvertFlightParams{Flight: flight, Legs: legs, PNRs: pnrs}), nil
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

		locationId, err := qtx.InsertLocation(ctx, sqlc.InsertLocationParams{
			Latitude:  airport.Location.Latitude,
			Longitude: airport.Location.Longitude,
		})
		if err != nil {
			return entity.Flight{}, err
		}
		err = qtx.InsertAirport(ctx, sqlc.InsertAirportParams{
			Iata:         airport.Iata,
			Name:         airport.Name,
			Municipality: airport.Municipality,
			LocationID:   &locationId,
		})
		if err != nil {
			return entity.Flight{}, err
		}
	}

	for _, leg := range flight.Legs {
		_, err := qtx.InsertFlightLeg(ctx, sqlc.InsertFlightLegParams{
			FlightID:          flightId,
			Origin:            leg.Origin.Iata,
			Destination:       leg.Destination.Iata,
			Airline:           leg.Airline,
			FlightNumber:      leg.FlightNumber,
			DepartureTime:     leg.DepartureDateTime,
			ArrivalTime:       leg.ArrivalDateTime,
			DurationInMinutes: leg.DurationInMinutes,
			Aircraft:          leg.Aircraft,
		})
		if err != nil {
			return entity.Flight{}, err
		}
	}

	for _, pnr := range flight.PNRs {
		_, err := qtx.InsertPNR(ctx, sqlc.InsertPNRParams{
			FlightID: flightId,
			Airline:  pnr.Airline,
			Pnr:      pnr.PNR,
		})
		if err != nil {
			return entity.Flight{}, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return entity.Flight{}, err
	}

	return r.GetFlightByID(ctx, flight.TripID, flightId)
}

func (r *FlightsRepo) DeleteFlight(ctx context.Context, tripID int32, flightID int32) error {
	return r.Queries.DeleteFlightByID(ctx, sqlc.DeleteFlightByIDParams{TripID: tripID, ID: flightID})
}
