package v1

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"kompass/internal/controller/http/v1/request"
	"kompass/internal/usecase"
	"kompass/pkg/logger"
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
		return errorResponse(ctx, fmt.Errorf("get trips: %w", err))
	}

	return ctx.Status(http.StatusOK).JSON(trips)
}

// @Summary     Get trip by ID
// @ID          getTrip
// @Tags  	    trips
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Success     200 {object} entity.Trip
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id} [get]
func (r *TripsV1) getTrip(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponse(ctx, fmt.Errorf("parse trip_id: %w", err))
	}

	trip, err := r.uc.GetTripByID(ctx.UserContext(), int32(tripID))
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusNotFound, fmt.Errorf("get trip [id=%d]: %w", tripID, err))
	}

	return ctx.Status(http.StatusOK).JSON(trip)
}

// @Summary     Add trip
// @ID          postTrip
// @Tags  	    trips
// @Accept      json
// @Produce     json
// @Param       request body request.Trip true "trip"
// @Success     204
// @Failure     500 {object} response.Error
// @Router      /trips [post]
func (r *TripsV1) postTrip(ctx *fiber.Ctx) error {
	var body request.Trip

	if err := ctx.BodyParser(&body); err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("parse request body: %w", err))
	}

	if err := r.v.Struct(body); err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("validate request body: %w", err))
	}

	_, err := r.uc.CreateTrip(ctx.UserContext(), body)
	if err != nil {
		return errorResponse(ctx, fmt.Errorf("create trip: %w", err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     Update trip
// @ID          putTrip
// @Tags  	    trips
// @Accept      json
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       request body request.Trip true "trip"
// @Success     204
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id} [put]
func (r *TripsV1) putTrip(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("parse trip_id: %w", err))
	}

	var body request.Trip

	if err := ctx.BodyParser(&body); err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("parse request body: %w", err))
	}

	if err := r.v.Struct(body); err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("validate request body: %w", err))
	}

	if err := r.uc.UpdateTrip(ctx.UserContext(), int32(tripID), body); err != nil {
		return errorResponse(ctx, fmt.Errorf("update trip with id %d: %w", tripID, err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     Delete trip
// @ID          deleteTrip
// @Tags  	    trips
// @Param       trip_id path int true "Trip ID"
// @Success     204
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id} [delete]
func (r *TripsV1) deleteTrip(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("parse trip_id: %w", err))
	}

	if err := r.uc.DeleteTrip(ctx.UserContext(), int32(tripID)); err != nil {
		return errorResponse(ctx, fmt.Errorf("delete trip with id %d: %w", tripID, err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}
