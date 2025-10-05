package accommodation

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"kompass/internal/controller/http/v1/request"
	"kompass/internal/entity"
	"kompass/internal/repo"
	"kompass/internal/usecase"
)

type UseCase struct {
	repo  repo.AccommodationRepo
	trips usecase.Trips
}

func New(r repo.AccommodationRepo, trips usecase.Trips) *UseCase {
	return &UseCase{
		repo:  r,
		trips: trips,
	}
}

func (uc *UseCase) GetAccommodationByID(ctx context.Context, tripID int32, accommodationID int32) (entity.Accommodation, error) {
	return uc.repo.GetAccommodationByID(ctx, tripID, accommodationID)
}

func (uc *UseCase) GetAllAccommodation(ctx context.Context, tripID int32) ([]entity.Accommodation, error) {
	return uc.repo.GetAllAccommodation(ctx, tripID)
}

func (uc *UseCase) CreateAccommodation(ctx context.Context, tripID int32, accommodation request.Accommodation) (entity.Accommodation, error) {
	if err := uc.trips.VerifyDatesInBounds(ctx, tripID, accommodation.DepartureDate, accommodation.ArrivalDate); err != nil {
		return entity.Accommodation{}, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return uc.repo.CreateAccommodation(ctx, entity.Accommodation{
		TripID:        tripID,
		Name:          accommodation.Name,
		ArrivalDate:   accommodation.ArrivalDate,
		DepartureDate: accommodation.DepartureDate,
		CheckInTime:   accommodation.CheckInTime,
		CheckOutTime:  accommodation.CheckOutTime,
		Description:   accommodation.Description,
		Address:       accommodation.Address,
		Location:      accommodation.Location,
		Price:         accommodation.Price,
	})
}

func (uc *UseCase) UpdateAccommodation(ctx context.Context, tripID int32, accommodationID int32, accommodation request.Accommodation) error {
	if err := uc.trips.VerifyDatesInBounds(ctx, tripID, accommodation.DepartureDate, accommodation.ArrivalDate); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return uc.repo.UpdateAccommodation(ctx, entity.Accommodation{
		ID:            accommodationID,
		TripID:        tripID,
		Name:          accommodation.Name,
		ArrivalDate:   accommodation.ArrivalDate,
		DepartureDate: accommodation.DepartureDate,
		CheckInTime:   accommodation.CheckInTime,
		CheckOutTime:  accommodation.CheckOutTime,
		Description:   accommodation.Description,
		Address:       accommodation.Address,
		Location:      accommodation.Location,
		Price:         accommodation.Price,
	})
}

func (uc *UseCase) DeleteAccommodation(ctx context.Context, tripID int32, accommodationID int32) error {
	return uc.repo.DeleteAccommodation(ctx, tripID, accommodationID)
}
