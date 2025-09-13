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

func (r *FlightsRepo) SaveFlightDetail(ctx context.Context, qtx *sqlc.Queries, transportationID int32, flight entity.FlightDetail) error {
	airportSet := make(map[string]entity.Airport)
	for _, leg := range flight.Legs {
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
