package postgres

import (
	"context"
	"fmt"
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

func (r *FlightsRepo) GetFlightDetail(ctx context.Context, transportationID int32) (entity.FlightDetail, error) {
	legs, err := r.Queries.GetFlightLegsByTransportationID(ctx, transportationID)
	if err != nil {
		return entity.FlightDetail{}, fmt.Errorf("get flight legs [t.id=%d] from db: %w", transportationID, err)
	}

	pnrs, err := r.Queries.GetPnrsByTransportationID(ctx, transportationID)
	if err != nil {
		return entity.FlightDetail{}, fmt.Errorf("get PNRs [t.id=%d] from db: %w", transportationID, err)
	}

	return entity.FlightDetail{
		Legs: r.c.ConvertFlightLegs(legs),
		PNRs: r.c.ConvertPnrs(pnrs),
	}, nil
}

func (r *FlightsRepo) CreateFlightDetail(ctx context.Context, qtx *sqlc.Queries, transportationID int32, flight entity.FlightDetail) error {
	err := r.saveAirports(ctx, qtx, flight.Legs)
	if err != nil {
		return fmt.Errorf("save airports: %w", err)
	}

	err = r.createFlightLegs(ctx, qtx, transportationID, flight)
	if err != nil {
		return fmt.Errorf("save flight legs: %w", err)
	}

	err = r.createPNRs(ctx, qtx, transportationID, flight)
	if err != nil {
		return fmt.Errorf("save pnrs: %w", err)
	}

	return nil
}

func (r *FlightsRepo) UpdateFlightLegs(ctx context.Context, flightLegs []entity.FlightLeg) error {
	tx, err := r.Db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := r.Queries.WithTx(tx)

	err = r.saveAirports(ctx, r.Queries, flightLegs)
	if err != nil {
		return fmt.Errorf("save airports: %w", err)
	}

	for _, leg := range flightLegs {
		err := qtx.UpdateFlightLeg(ctx, sqlc.UpdateFlightLegParams{
			ID:                leg.ID,
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
			return fmt.Errorf("update flight leg [id=%d]: %w", leg.ID, err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (r *FlightsRepo) saveAirports(ctx context.Context, qtx *sqlc.Queries, flightLegs []entity.FlightLeg) error {
	airportSet := make(map[string]entity.Airport)
	for _, leg := range flightLegs {
		airportSet[leg.Origin.Iata] = leg.Origin
		airportSet[leg.Destination.Iata] = leg.Destination
	}

	for iata := range airportSet {
		airport := airportSet[iata]

		locationId, err := SaveLocation(ctx, qtx, airport.Location)
		if err != nil {
			return fmt.Errorf("save location: %w", err)
		}
		err = qtx.InsertAirport(ctx, sqlc.InsertAirportParams{
			Iata:         airport.Iata,
			Name:         airport.Name,
			Municipality: airport.Municipality,
			LocationID:   &locationId,
		})
		if err != nil {
			return fmt.Errorf("insert airport: %w", err)
		}
	}

	return nil
}

func (r *FlightsRepo) createFlightLegs(ctx context.Context, qtx *sqlc.Queries, transportationID int32, flight entity.FlightDetail) error {
	for _, leg := range flight.Legs {
		_, err := qtx.InsertFlightLeg(ctx, sqlc.InsertFlightLegParams{
			TransportationID:  transportationID,
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
			return fmt.Errorf("insert leg: %w", err)
		}
	}

	return nil
}

func (r *FlightsRepo) createPNRs(ctx context.Context, qtx *sqlc.Queries, transportationID int32, flight entity.FlightDetail) error {
	for _, pnr := range flight.PNRs {
		_, err := qtx.InsertPNR(ctx, sqlc.InsertPNRParams{
			TransportationID: transportationID,
			Airline:          pnr.Airline,
			Pnr:              pnr.PNR,
		})
		if err != nil {
			return fmt.Errorf("insert pnr: %w", err)
		}
	}

	return nil
}
