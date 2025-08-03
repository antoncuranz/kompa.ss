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

func (uc *UseCase) GetAccommodationByID(ctx context.Context, id int32) (entity.Accommodation, error) {
	return uc.repo.GetAccommodationByID(ctx, id)
}

func (uc *UseCase) GetAllAccommodation(ctx context.Context) ([]entity.Accommodation, error) {
	return uc.repo.GetAllAccommodation(ctx)
}

func (uc *UseCase) CreateAccommodation(ctx context.Context, accommodation request.Accommodation) (entity.Accommodation, error) {
	return uc.repo.SaveAccommodation(ctx, entity.Accommodation{
		TripID:        accommodation.TripID,
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
