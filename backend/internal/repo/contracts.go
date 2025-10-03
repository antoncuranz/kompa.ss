// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"cloud.google.com/go/civil"
	"context"
	"github.com/google/uuid"
	"github.com/paulmach/orb/geojson"
	"kompass/internal/controller/http/v1/request"
	"kompass/internal/entity"
	"kompass/pkg/sqlc"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	UserRepo interface {
		GetUsers(ctx context.Context) ([]entity.User, error)
		GetUserByID(ctx context.Context, id int32) (entity.User, error)
		GetUserByJwtSub(ctx context.Context, sub uuid.UUID) (entity.User, error)
		CreateUser(ctx context.Context, user entity.User) (entity.User, error)
		HasReadPermission(ctx context.Context, userID, tripID int32) (bool, error)
		HasWritePermission(ctx context.Context, userID, tripID int32) (bool, error)
		IsTripOwner(ctx context.Context, userID int32, tripID int32) (bool, error)
	}

	TripsRepo interface {
		GetTrips(ctx context.Context, userID int32) ([]entity.Trip, error)
		GetTripByID(ctx context.Context, id int32) (entity.Trip, error)
		CreateTrip(ctx context.Context, trip entity.Trip) (entity.Trip, error)
		UpdateTrip(ctx context.Context, trip entity.Trip) error
		DeleteTrip(ctx context.Context, tripID int32) error
	}

	TransportationRepo interface {
		GetAllTransportation(ctx context.Context, tripID int32) ([]entity.Transportation, error)
		GetTransportationByID(ctx context.Context, tripID int32, id int32) (entity.Transportation, error)
		SaveTransportation(ctx context.Context, transportation entity.Transportation) (entity.Transportation, error)
		DeleteTransportation(ctx context.Context, tripID int32, flightID int32) error
		GetAllGeoJson(ctx context.Context, tripID int32) ([]geojson.FeatureCollection, error)
		SaveGeoJson(ctx context.Context, transportationID int32, geoJson *geojson.FeatureCollection) error
	}

	FlightsRepo interface {
		GetFlightDetail(ctx context.Context, transportationID int32) (entity.FlightDetail, error)
		SaveFlightDetail(ctx context.Context, qtx *sqlc.Queries, transportationID int32, flight entity.FlightDetail) error
		UpdateFlightLegs(ctx context.Context, flightLegs []entity.FlightLeg) error
	}

	ActivitiesRepo interface {
		GetActivities(ctx context.Context, tripID int32) ([]entity.Activity, error)
		GetActivityByID(ctx context.Context, tripID int32, activityID int32) (entity.Activity, error)
		SaveActivity(ctx context.Context, activity entity.Activity) (entity.Activity, error)
		UpdateActivity(ctx context.Context, activity entity.Activity) error
		DeleteActivity(ctx context.Context, tripID int32, activityID int32) error
	}

	AccommodationRepo interface {
		GetAllAccommodation(ctx context.Context, tripID int32) ([]entity.Accommodation, error)
		GetAccommodationByID(ctx context.Context, tripID int32, id int32) (entity.Accommodation, error)
		SaveAccommodation(ctx context.Context, accommodation entity.Accommodation) (entity.Accommodation, error)
		UpdateAccommodation(ctx context.Context, accommodation entity.Accommodation) error
		DeleteAccommodation(ctx context.Context, tripID int32, accommodationID int32) error
	}

	AttachmentsRepo interface {
		GetAttachments(ctx context.Context, tripID int32) ([]entity.Attachment, error)
		GetAttachmentByID(ctx context.Context, tripID int32, attachmentID int32) (entity.Attachment, error)
		SaveAttachment(ctx context.Context, attachment entity.Attachment) (entity.Attachment, error)
		DeleteAttachment(ctx context.Context, tripID int32, attachmentID int32) error
	}

	AerodataboxWebAPI interface {
		RetrieveFlightLeg(ctx context.Context, date civil.Date, flightNumber string, origin *string) (entity.FlightLeg, error)
	}

	DbVendoWebAPI interface {
		RetrieveLocation(ctx context.Context, query string) (entity.TrainStation, error)
		RetrieveJourney(ctx context.Context, journey request.TrainJourney) (entity.TrainDetail, error)
		RetrievePolylines(ctx context.Context, refreshToken string) ([]geojson.FeatureCollection, error)
	}
)
