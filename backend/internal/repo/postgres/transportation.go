package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/paulmach/orb/geojson"
	"kompass/internal/entity"
	"kompass/internal/repo/postgres/converter"
	"kompass/pkg/postgres"
	"kompass/pkg/sqlc"
	"net/http"

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
		transportation, err := r.convertTransportationRow(ctx, row.Transportation, row.Location, row.Location_2)
		if err != nil {
			return []entity.Transportation{}, fmt.Errorf("convert row to transportation: %w", err)
		}

		result = append(result, transportation)
	}
	return result, nil
}

func (r *TransportationRepo) GetTransportationByID(ctx context.Context, tripID int32, transportationID int32) (entity.Transportation, error) {
	row, err := r.Queries.GetTransportationByID(ctx, sqlc.GetTransportationByIDParams{TripID: tripID, ID: transportationID})
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return entity.Transportation{}, fiber.NewError(http.StatusNotFound, "transportation not found")
		}
		return entity.Transportation{}, fmt.Errorf("get transportation [id=%d] from db: %w", transportationID, err)
	}

	transportation, err := r.convertTransportationRow(ctx, row.Transportation, row.Location, row.Location_2)
	if err != nil {
		return entity.Transportation{}, fmt.Errorf("convert row to transportation: %w", err)
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

	transportationID, err := r.saveTransportation(ctx, qtx, transportation, originId, destinationId)
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return entity.Transportation{}, fiber.NewError(http.StatusNotFound, "transportation not found")
		}
		return entity.Transportation{}, fmt.Errorf("save transportation: %w", err)
	}

	if transportation.GenericDetail != nil {
		if err := r.saveGenericDetail(ctx, qtx, transportationID, *transportation.GenericDetail); err != nil {
			return entity.Transportation{}, fmt.Errorf("save generic detail: %w", err)
		}
	}

	// it is not possible to update flight or train detail via the transportation api yet
	if transportation.ID == 0 {
		if transportation.FlightDetail != nil {
			if err := r.flights.CreateFlightDetail(ctx, qtx, transportationID, *transportation.FlightDetail); err != nil {
				return entity.Transportation{}, fmt.Errorf("save flight detail: %w", err)
			}
		}
		if transportation.TrainDetail != nil {
			if err := r.trains.CreateTrainDetail(ctx, qtx, transportationID, *transportation.TrainDetail); err != nil {
				return entity.Transportation{}, fmt.Errorf("save train detail: %w", err)
			}
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return entity.Transportation{}, fmt.Errorf("commit tx: %w", err)
	}
	return r.GetTransportationByID(ctx, transportation.TripID, transportationID)
}

func (r *TransportationRepo) DeleteTransportation(ctx context.Context, tripID int32, transportationID int32) error {
	_, err := r.Queries.DeleteTransportationByID(ctx, sqlc.DeleteTransportationByIDParams{TripID: tripID, ID: transportationID})
	if err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return fiber.NewError(http.StatusNotFound, "transportation not found")
		}
		return fmt.Errorf("delete transportation [id=%d]: %w", transportationID, err)
	}

	return nil
}

func (r *TransportationRepo) GetAllGeoJson(ctx context.Context, tripID int32) ([]geojson.FeatureCollection, error) {
	rows, err := r.Queries.GetAllGeoJson(ctx, tripID)
	if err != nil {
		return nil, fmt.Errorf("get all geojson from db [t.id=%d]: %w", tripID, err)
	}

	allGeoJson := []geojson.FeatureCollection{}
	for _, bytes := range rows {
		geoJson, err := geojson.UnmarshalFeatureCollection(bytes)
		if err != nil {
			return nil, fmt.Errorf("unmarshall GeoJson: %w", err)
		}
		allGeoJson = append(allGeoJson, *geoJson)
	}
	return allGeoJson, nil
}

func (r *TransportationRepo) SaveGeoJson(ctx context.Context, transportationID int32, geoJson *geojson.FeatureCollection) error {
	marshalledGeoJson, err := geoJson.MarshalJSON()
	if err != nil {
		return fmt.Errorf("marshal geojson: %w", err)
	}

	return r.Queries.UpsertGeoJson(ctx, sqlc.UpsertGeoJsonParams{
		TransportationID: transportationID,
		Geojson:          marshalledGeoJson,
	})
}

func (r *TransportationRepo) saveTransportation(ctx context.Context, queries *sqlc.Queries, transportation entity.Transportation, originID int32, destinationID int32) (int32, error) {
	if transportation.ID != 0 {
		return transportation.ID, queries.UpdateTransportation(ctx, sqlc.UpdateTransportationParams{
			ID:            transportation.ID,
			Type:          transportation.Type.String(),
			OriginID:      originID,
			DestinationID: destinationID,
			DepartureTime: transportation.DepartureDateTime,
			ArrivalTime:   transportation.ArrivalDateTime,
			Price:         transportation.Price,
		})
	} else {
		return queries.InsertTransportation(ctx, sqlc.InsertTransportationParams{
			TripID:        transportation.TripID,
			Type:          transportation.Type.String(),
			OriginID:      originID,
			DestinationID: destinationID,
			DepartureTime: transportation.DepartureDateTime,
			ArrivalTime:   transportation.ArrivalDateTime,
			Price:         transportation.Price,
		})
	}
}

func (r *TransportationRepo) convertTransportationRow(ctx context.Context, transportation sqlc.Transportation, origin sqlc.Location, destination sqlc.Location) (entity.Transportation, error) {
	converted := r.c.ConvertTransportation(converter.ConvertTransportationParams{
		Transportation: transportation,
		Origin:         origin,
		Destination:    destination,
	})

	if converted.Type == entity.FLIGHT {
		flightDetail, err := r.flights.GetFlightDetail(ctx, converted.ID)
		if err != nil {
			return entity.Transportation{}, fmt.Errorf("get flightDetail [t.id=%d]: %w", converted.ID, err)
		}
		converted.FlightDetail = &flightDetail
	} else if converted.Type == entity.TRAIN {
		trainDetail, err := r.trains.GetTrainDetail(ctx, converted.ID)
		if err != nil {
			return entity.Transportation{}, fmt.Errorf("get trainDetail [t.id=%d]: %w", converted.ID, err)
		}
		converted.TrainDetail = &trainDetail
	} else {
		genericDetail, err := r.getGenericDetail(ctx, converted.ID)
		if err != nil {
			return entity.Transportation{}, fmt.Errorf("get genericDetail [t.id=%d]: %w", converted.ID, err)
		}
		converted.GenericDetail = &genericDetail
	}

	return converted, nil
}

func (r *TransportationRepo) saveGenericDetail(ctx context.Context, qtx *sqlc.Queries, transportationID int32, detail entity.GenericDetail) error {
	return qtx.UpsertGenericTransportationDetail(ctx, sqlc.UpsertGenericTransportationDetailParams{
		TransportationID:   transportationID,
		Name:               detail.Name,
		OriginAddress:      detail.OriginAddress,
		DestinationAddress: detail.DestinationAddress,
	})
}

func (r *TransportationRepo) getGenericDetail(ctx context.Context, transportationID int32) (entity.GenericDetail, error) {
	row, err := r.Queries.GetGenericDetailByTransportationID(ctx, transportationID)
	if err != nil {
		return entity.GenericDetail{}, fmt.Errorf("get generic detail: %w", err)
	}

	return entity.GenericDetail{
		Name:               row.Name,
		OriginAddress:      row.OriginAddress,
		DestinationAddress: row.DestinationAddress,
	}, err
}
