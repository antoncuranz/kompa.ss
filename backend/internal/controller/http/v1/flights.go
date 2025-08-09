package v1

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"travel-planner/internal/controller/http/v1/request"
	"travel-planner/internal/usecase"
	"travel-planner/pkg/logger"
)

type FlightsV1 struct {
	uc  usecase.Flights
	log logger.Interface
	v   *validator.Validate
}

// @Summary     Get all flights
// @ID          getFlights
// @Tags  	    flights
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Success     200 {object} []entity.Flight
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/flights [get]
func (r *FlightsV1) getFlights(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "unable to parse trip_id")
	}

	flights, err := r.uc.GetFlights(ctx.UserContext(), int32(tripID))
	if err != nil {
		r.log.Error(err, "http - v1 - getFlights")
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(flights)
}

// @Summary     Get flight by ID
// @ID          getFlight
// @Tags  	    flights
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       flight_id path int true "Flight ID"
// @Success     200 {object} entity.Flight
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/flights/{flight_id} [get]
func (r *FlightsV1) getFlight(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "unable to parse trip_id")
	}
	flightID, err := ctx.ParamsInt("flight_id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "unable to parse flight_id")
	}

	flight, err := r.uc.GetFlightByID(ctx.UserContext(), int32(tripID), int32(flightID))
	if err != nil {
		return errorResponse(ctx, http.StatusNotFound, fmt.Sprintf("unable to find flight with id %d", flightID))
	}

	return ctx.Status(http.StatusOK).JSON(flight)
}

// @Summary     Add flight
// @ID          postFlight
// @Tags  	    flights
// @Accept      json
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       request body request.Flight true "flight"
// @Success     200 {object} entity.Flight
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/flights [post]
func (r *FlightsV1) postFlight(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "unable to parse trip_id")
	}

	var body request.Flight

	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		fmt.Println(err)
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	_, err = r.uc.CreateFlight(ctx.UserContext(), int32(tripID), body)
	if err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusOK)
}
