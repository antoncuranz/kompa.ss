package v1

import (
	"backplate/internal/usecase"
	"backplate/pkg/logger"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
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
// @Success     200 {object} []entity.Flight
// @Failure     500 {object} response.Error
// @Router      /flights [get]
func (r *FlightsV1) getFlights(ctx *fiber.Ctx) error {
	flights, err := r.uc.GetFlights(ctx.UserContext())
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
// @Param       flight_id path int true "Flight ID"
// @Success     200 {object} entity.Flight
// @Failure     500 {object} response.Error
// @Router      /flights/{flight_id} [get]
func (r *FlightsV1) getFlight(ctx *fiber.Ctx) error {
	flightID, err := ctx.ParamsInt("flight_id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "unable to parse flight_id")
	}

	flight, err := r.uc.GetFlightByID(ctx.UserContext(), int32(flightID))
	if err != nil {
		return errorResponse(ctx, http.StatusNotFound, fmt.Sprintf("unable to find flight with id %d", flightID))
	}

	return ctx.Status(http.StatusOK).JSON(flight)
}
