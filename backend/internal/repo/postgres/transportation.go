package postgres

import (
	"context"
	"fmt"
	"github.com/paulmach/orb/geojson"
	"kompass/internal/entity"
	"kompass/internal/repo/postgres/converter"
	"kompass/pkg/postgres"
	"kompass/pkg/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransportationRepo struct {
	Db      *pgxpool.Pool
	Queries *sqlc.Queries
	c       converter.TransportationConverter
	flights *FlightsRepo
	trains  *TrainsRepo
}

func NewTransportationRepo(pg *postgres.Postgres, flights *FlightsRepo, trains *TrainsRepo) *TransportationRepo {
	return &TransportationRepo{
		pg.Pool,
		sqlc.New(pg.Pool),
		&converter.TransportationConverterImpl{},
		flights,
		trains,
	}
}

func (r *TransportationRepo) GetAllTransportation(ctx context.Context, tripID int32) ([]entity.Transportation, error) {
	rows, err := r.Queries.GetAllTransportation(ctx, tripID)
	if err != nil {
		return []entity.Transportation{}, fmt.Errorf("get all transportation from db: %w", err)
	}

	result := []entity.Transportation{}
	for _, row := range rows {
		transportation := r.c.ConvertTransportation(converter.ConvertTransportationParams{
			Transportation: row.Transportation,
			Origin:         row.Location,
			Destination:    row.Location_2,
		})

		if transportation.Type == entity.PLANE {
			flightDetail, err := r.flights.GetFlightDetail(ctx, transportation.ID)
			if err != nil {
				return []entity.Transportation{}, fmt.Errorf("get flightDetail [t.id=%d]: %w", transportation.ID, err)
			}
			transportation.FlightDetail = &flightDetail
		} else if transportation.Type == entity.TRAIN {
			trainDetail, err := r.trains.GetTrainDetail(ctx, transportation.ID)
			if err != nil {
				return []entity.Transportation{}, fmt.Errorf("get trainDetail [t.id=%d]: %w", transportation.ID, err)
			}
			transportation.TrainDetail = &trainDetail
		}

		result = append(result, transportation)
	}
	return result, nil
}

func (r *TransportationRepo) GetTransportationByID(ctx context.Context, tripID int32, transportationID int32) (entity.Transportation, error) {
	row, err := r.Queries.GetTransportationByID(ctx, sqlc.GetTransportationByIDParams{TripID: tripID, ID: transportationID})
	if err != nil {
		return entity.Transportation{}, fmt.Errorf("get transportation [id=%d] from db: %w", transportationID, err)
	}

	transportation := r.c.ConvertTransportation(converter.ConvertTransportationParams{
		Transportation: row.Transportation,
		Origin:         row.Location,
		Destination:    row.Location_2,
	})

	if transportation.Type == entity.PLANE {
		flightDetail, err := r.flights.GetFlightDetail(ctx, transportation.ID)
		if err != nil {
			return entity.Transportation{}, fmt.Errorf("get flightDetail [t.id=%d]: %w", transportation.ID, err)
		}
		transportation.FlightDetail = &flightDetail
	} else if transportation.Type == entity.TRAIN {
		trainDetail, err := r.trains.GetTrainDetail(ctx, transportation.ID)
		if err != nil {
			return entity.Transportation{}, fmt.Errorf("get trainDetail [t.id=%d]: %w", transportation.ID, err)
		}
		transportation.TrainDetail = &trainDetail
	}

	return transportation, nil
}

func (r *TransportationRepo) SaveTransportation(ctx context.Context, transportation entity.Transportation) (entity.Transportation, error) {

	tx, err := r.Db.Begin(ctx)
	if err != nil {
		return entity.Transportation{}, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := r.Queries.WithTx(tx)

	originId, err := SaveLocation(ctx, qtx, transportation.Origin)
	if err != nil {
		return entity.Transportation{}, fmt.Errorf("save origin location: %w", err)
	}
	destinationId, err := SaveLocation(ctx, qtx, transportation.Destination)
	if err != nil {
		return entity.Transportation{}, fmt.Errorf("save destination location: %w", err)
	}

	transportationID, err := qtx.InsertTransportation(ctx, sqlc.InsertTransportationParams{
		TripID:        transportation.TripID,
		Type:          transportation.Type.String(),
		OriginID:      originId,
		DestinationID: destinationId,
		DepartureTime: transportation.DepartureDateTime,
		ArrivalTime:   transportation.ArrivalDateTime,
		Price:         transportation.Price,
	})
	if err != nil {
		return entity.Transportation{}, fmt.Errorf("insert transportation: %w", err)
	}

	if transportation.FlightDetail != nil {
		if err := r.flights.SaveFlightDetail(ctx, qtx, transportationID, *transportation.FlightDetail); err != nil {
			return entity.Transportation{}, fmt.Errorf("save flight detail: %w", err)
		}
	}
	if transportation.TrainDetail != nil {
		if err := r.trains.SaveTrainDetail(ctx, qtx, transportationID, *transportation.TrainDetail); err != nil {
			return entity.Transportation{}, fmt.Errorf("save train detail: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return entity.Transportation{}, fmt.Errorf("commit tx: %w", err)
	}

	return r.GetTransportationByID(ctx, transportation.TripID, transportationID)
}

func (r *TransportationRepo) SaveGeoJson(ctx context.Context, transportationID int32, geoJson *geojson.FeatureCollection) error {
	marshalledGeoJson, err := geoJson.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal geojson: %w", err)
	}

	return r.Queries.InsertGeoJson(ctx, sqlc.InsertGeoJsonParams{
		TransportationID: transportationID,
		Geojson:          marshalledGeoJson,
	})
}

func (r *TransportationRepo) DeleteTransportation(ctx context.Context, tripID int32, transportationID int32) error {
	return r.Queries.DeleteTransportationByID(ctx, sqlc.DeleteTransportationByIDParams{TripID: tripID, ID: transportationID})
}
