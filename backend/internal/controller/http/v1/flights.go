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

type FlightsV1 struct {
	uc  usecase.Flights
	log logger.Interface
	v   *validator.Validate
}

// @Summary     Add flight
// @ID          postFlight
// @Tags  	    flights
// @Accept      json
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       request body request.Flight true "flight"
// @Success     204
// @Failure     403 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id}/flights [post]
func (r *FlightsV1) postFlight(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "unable to parse trip_id")
	}

	body, err := ParseAndValidateRequestBody[request.Flight](ctx, r.v)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid request body")
	}

	_, err = r.uc.CreateFlight(ctx.UserContext(), int32(tripID), *body)
	if err != nil {
		return fmt.Errorf("create flight: %w", err)
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     Update flight
// @ID          putFlight
// @Tags  	    flights
// @Param       trip_id path int true "Trip ID"
// @Param       flight_id path int true "Flight ID"
// @Success     204
// @Failure     403 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id}/flights/{flight_id} [put]
func (r *FlightsV1) putFlight(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "unable to parse trip_id")
	}
	flightID, err := ctx.ParamsInt("flight_id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "unable to parse flight_id")
	}

	if err := r.uc.UpdateFlight(ctx.UserContext(), int32(tripID), int32(flightID)); err != nil {
		return fmt.Errorf("update flight with id %d: %w", flightID, err)
	}

	return ctx.SendStatus(http.StatusNoContent)
}
