package v1

import (
	"fmt"
	"kompass/internal/controller/http/v1/request"
	"kompass/internal/usecase"
	"kompass/pkg/logger"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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
// @Security    bearerauth
// @Router      /trips [get]
func (r *TripsV1) getTrips(ctx *fiber.Ctx) error {
	trips, err := r.uc.GetTrips(ctx.UserContext(), userIdFromCtx(ctx))
	if err != nil {
		return fmt.Errorf("get trips: %w", err)
	}

	return ctx.Status(http.StatusOK).JSON(trips)
}

// @Summary     Get trip by ID
// @ID          getTrip
// @Tags  	    trips
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Success     200 {object} entity.Trip
// @Failure     400 {object} response.Error
// @Failure     403 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id} [get]
func (r *TripsV1) getTrip(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "parse trip_id")
	}

	trip, err := r.uc.GetTripByID(ctx.UserContext(), int32(tripID))
	if err != nil { // TODO: differ between 404 and 500
		return fmt.Errorf("get trip [id=%d]: %w", tripID, err)
		//return ErrorResponseWithStatus(ctx, http.StatusNotFound, fmt.Errorf("get trip [id=%d]: %w", tripID, err))
	}

	return ctx.Status(http.StatusOK).JSON(trip)
}

// @Summary     Add trip
// @ID          postTrip
// @Tags  	    trips
// @Accept      json
// @Produce     json
// @Param       request body request.Trip true "trip"
// @Success     200 {object} entity.Trip
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips [post]
func (r *TripsV1) postTrip(ctx *fiber.Ctx) error {
	body, err := ParseAndValidateRequestBody[request.Trip](ctx, r.v)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "parse request body")
	}

	trip, err := r.uc.CreateTrip(ctx.UserContext(), userIdFromCtx(ctx), *body)
	if err != nil {
		return fmt.Errorf("create trip: %w", err)
	}

	return ctx.Status(http.StatusOK).JSON(trip)
}

// @Summary     Update trip
// @ID          putTrip
// @Tags  	    trips
// @Accept      json
// @Param       trip_id path int true "Trip ID"
// @Param       request body request.Trip true "trip"
// @Success     204
// @Failure     403 {object} response.Error
// @Failure     400 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id} [put]
func (r *TripsV1) putTrip(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "parse trip_id")
	}

	body, err := ParseAndValidateRequestBody[request.Trip](ctx, r.v)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid request body")
	}

	if err := r.uc.UpdateTrip(ctx.UserContext(), int32(tripID), *body); err != nil { // TODO: 404
		return fmt.Errorf("update trip with id %d: %w", tripID, err)
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     Delete trip
// @ID          deleteTrip
// @Tags  	    trips
// @Param       trip_id path int true "Trip ID"
// @Success     204
// @Failure     400 {object} response.Error
// @Failure     403 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id} [delete]
func (r *TripsV1) deleteTrip(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "parse trip_id")
	}

	if err := r.uc.DeleteTrip(ctx.UserContext(), int32(tripID)); err != nil { // TODO: 404
		return fmt.Errorf("delete trip with id %d: %w", tripID, err)
	}

	return ctx.SendStatus(http.StatusNoContent)
}
