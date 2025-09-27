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

type TrainsV1 struct {
	uc  usecase.Trains
	log logger.Interface
	v   *validator.Validate
}

// @Summary     Get train station
// @ID          getTrainStation
// @Tags  	    trains
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       query query string true "station query"
// @Success     200 {object} entity.TrainStation
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id}/trains/stations [get]
func (r *TrainsV1) getTrainStation(ctx *fiber.Ctx) error {
	query := ctx.Query("query")
	location, err := r.uc.RetrieveLocation(ctx.Context(), query)
	if err != nil {
		return ErrorResponse(ctx, fmt.Errorf("retrieve location: %w", err))
	}

	return ctx.Status(http.StatusOK).JSON(location)
}

// @Summary     Add train journey
// @ID          postTrainJourney
// @Tags  	    trains
// @Accept      json
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       request body request.TrainJourney true "train journey"
// @Success     200 {object} entity.Transportation
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id}/trains [post]
func (r *TrainsV1) postTrainJourney(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return ErrorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse trip_id: %w", err))
	}

	body, err := ParseAndValidateRequestBody[request.TrainJourney](ctx, r.v)
	if err != nil {
		return ErrorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("parse request body: %w", err))
	}

	transportation, err := r.uc.CreateTrainJourney(ctx.Context(), userIdFromCtx(ctx), int32(tripID), *body)
	if err != nil {
		return ErrorResponse(ctx, fmt.Errorf("retrieve journey: %w", err))
	}

	return ctx.Status(http.StatusOK).JSON(transportation)
}
