package v1

import (
	"backplate/internal/usecase"
	"backplate/pkg/logger"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type TripsV1 struct {
	uc  usecase.Trips
	log logger.Interface
	v   *validator.Validate
}

// @Summary     Get all trips
// @ID          getTrips
// @Tags  	    trips
// @Produce     json
// @Success     200 {object} []entity.Trip
// @Failure     500 {object} response.Error
// @Router      /trips [get]
func (r *TripsV1) getTrips(ctx *fiber.Ctx) error {
	trips, err := r.uc.GetTrips(ctx.UserContext())
	if err != nil {
		r.log.Error(err, "http - v1 - getTrips")
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(trips)
}

// @Summary     Get trip by ID
// @ID          getTrip
// @Tags  	    trips
// @Produce     json
// @Param       trip_id path string true "Trip ID"
// @Success     200 {object} entity.Trip
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id} [get]
func (r *TripsV1) getTrip(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "unable to parse trip_id")
	}

	trip, err := r.uc.GetTripByID(ctx.UserContext(), int32(tripID))
	if err != nil {
		return errorResponse(ctx, http.StatusNotFound, fmt.Sprintf("unable to find trip with id %d", tripID))
	}

	return ctx.Status(http.StatusOK).JSON(trip)
}
