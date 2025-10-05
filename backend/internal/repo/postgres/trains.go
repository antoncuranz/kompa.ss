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

type TrainsRepo struct {
	Db      *pgxpool.Pool
	Queries *sqlc.Queries
	c       converter.TrainConverter
}

func NewTrainsRepo(pg *postgres.Postgres) *TrainsRepo {
	return &TrainsRepo{
		pg.Pool,
		sqlc.New(pg.Pool),
		&converter.TrainConverterImpl{},
	}
}

func (r *TrainsRepo) GetTrainDetail(ctx context.Context, transportationID int32) (entity.TrainDetail, error) {
	detail, err := r.Queries.GetTrainDetailByTransportationID(ctx, transportationID)
	if err != nil {
		return entity.TrainDetail{}, fmt.Errorf("get train detail [t.id=%d] from db: %w", transportationID, err)
	}

	legs, err := r.Queries.GetTrainLegsByTransportationID(ctx, transportationID)
	if err != nil {
		return entity.TrainDetail{}, fmt.Errorf("get train legs [t.id=%d] from db: %w", transportationID, err)
	}

	return entity.TrainDetail{
		RefreshToken: detail.RefreshToken,
		Legs:         r.c.ConvertTrainLegs(legs),
	}, nil
}

func (r *TrainsRepo) CreateTrainDetail(ctx context.Context, qtx *sqlc.Queries, transportationID int32, train entity.TrainDetail) error {
	err := qtx.InsertTrainDetail(ctx, sqlc.InsertTrainDetailParams{
		TransportationID: transportationID,
		RefreshToken:     train.RefreshToken,
	})
	if err != nil {
		return fmt.Errorf("save train detail: %w", err)
	}

	trainStationSet := make(map[string]entity.TrainStation)
	for _, leg := range train.Legs {
		trainStationSet[leg.Origin.ID] = leg.Origin
		trainStationSet[leg.Destination.ID] = leg.Destination
	}

	for id := range trainStationSet {
		trainStation := trainStationSet[id]

		exists, err := qtx.TrainStationExists(ctx, id)
		if err != nil {
			return fmt.Errorf("check trainStation exists: %w", err)
		}
		if exists {
			continue
		}

		locationId, err := SaveLocation(ctx, qtx, trainStation.Location)
		if err != nil {
			return fmt.Errorf("save location: %w", err)
		}
		err = qtx.InsertTrainStation(ctx, sqlc.InsertTrainStationParams{
			ID:         trainStation.ID,
			Name:       trainStation.Name,
			LocationID: &locationId,
		})
		if err != nil {
			return fmt.Errorf("insert train station: %w", err)
		}
	}

	for _, leg := range train.Legs {
		_, err := qtx.InsertTrainLeg(ctx, sqlc.InsertTrainLegParams{
			TransportationID:  transportationID,
			Origin:            leg.Origin.ID,
			Destination:       leg.Destination.ID,
			DepartureTime:     leg.DepartureDateTime,
			ArrivalTime:       leg.ArrivalDateTime,
			DurationInMinutes: leg.DurationInMinutes,
			LineName:          leg.LineName,
			OperatorName:      leg.OperatorName,
		})
		if err != nil {
			return fmt.Errorf("insert leg: %w", err)
		}
	}

	return nil
}
