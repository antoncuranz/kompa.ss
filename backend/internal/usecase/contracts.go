// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
	"kompass/internal/controller/http/v1/request"
	"kompass/internal/entity"
	"mime/multipart"
)

//go:generate mockgen -source=contracts.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	UseCases struct {
		Users          Users
		Trips          Trips
		Flights        Flights
		Transportation Transportation
		Activities     Activities
		Accommodation  Accommodation
		Attachments    Attachments
	}

	Users interface {
		GetUsers(ctx context.Context) ([]entity.User, error)
		GetUserByID(ctx context.Context, id int32) (entity.User, error)
	}

	Trips interface {
		GetTrips(ctx context.Context) ([]entity.Trip, error)
		GetTripByID(ctx context.Context, id int32) (entity.Trip, error)
		CreateTrip(ctx context.Context, trip request.Trip) (entity.Trip, error)
		UpdateTrip(ctx context.Context, tripID int32, trip request.Trip) error
		DeleteTrip(ctx context.Context, tripID int32) error
	}

	Transportation interface {
		GetAllTransportation(ctx context.Context, tripID int32) ([]entity.Transportation, error)
		GetTransportationByID(ctx context.Context, tripID int32, transportationID int32) (entity.Transportation, error)
		DeleteTransportation(ctx context.Context, tripID int32, transportationID int32) error
	}

	Flights interface {
		CreateFlight(ctx context.Context, tripID int32, flight request.Flight) (entity.Transportation, error)
		UpdateFlight(ctx context.Context, tripID int32, flightID int32, flight request.Flight) error
	}

	Activities interface {
		GetActivities(ctx context.Context, tripID int32) ([]entity.Activity, error)
		GetActivityByID(ctx context.Context, tripID int32, activityID int32) (entity.Activity, error)
		CreateActivity(ctx context.Context, tripID int32, activity request.Activity) (entity.Activity, error)
		UpdateActivity(ctx context.Context, tripID int32, activityID int32, activity request.Activity) error
		DeleteActivity(ctx context.Context, tripID int32, activityID int32) error
	}

	Accommodation interface {
		GetAllAccommodation(ctx context.Context, tripID int32) ([]entity.Accommodation, error)
		GetAccommodationByID(ctx context.Context, tripID int32, id int32) (entity.Accommodation, error)
		CreateAccommodation(ctx context.Context, tripID int32, accommodation request.Accommodation) (entity.Accommodation, error)
		UpdateAccommodation(ctx context.Context, tripID int32, accommodationID int32, accommodation request.Accommodation) error
		DeleteAccommodation(ctx context.Context, tripID int32, accommodationID int32) error
	}

	Attachments interface {
		GetAttachments(ctx context.Context, tripID int32) ([]entity.Attachment, error)
		GetAttachmentByID(ctx context.Context, tripID int32, attachmentID int32) (entity.Attachment, error)
		CreateAttachment(ctx context.Context, tripID int32, attachment *multipart.FileHeader) (entity.Attachment, error)
		DeleteAttachment(ctx context.Context, tripID int32, attachmentID int32) error
	}
)
