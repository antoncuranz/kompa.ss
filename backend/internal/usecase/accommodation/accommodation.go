package accommodation

import (
	"context"
	"travel-planner/internal/controller/http/v1/request"
	"travel-planner/internal/entity"
	"travel-planner/internal/repo"
)

type UseCase struct {
	repo repo.AccommodationRepo
}

func New(r repo.AccommodationRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) GetAccommodationByID(ctx context.Context, tripID int32, accommodationID int32) (entity.Accommodation, error) {
	return uc.repo.GetAccommodationByID(ctx, tripID, accommodationID)
}

func (uc *UseCase) GetAllAccommodation(ctx context.Context, tripID int32) ([]entity.Accommodation, error) {
	return uc.repo.GetAllAccommodation(ctx, tripID)
}

func (uc *UseCase) CreateAccommodation(ctx context.Context, tripID int32, accommodation request.Accommodation) (entity.Accommodation, error) {
	return uc.repo.SaveAccommodation(ctx, entity.Accommodation{
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
