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

// @Summary     Add train journey
// @ID          postTrainJourney
// @Tags  	    trains
// @Accept      json
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       request body request.TrainJourney true "train journey"
// @Success     200 {object} entity.Transportation
// @Failure     403 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id}/trains [post]
func (r *TrainsV1) postTrainJourney(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "unable to parse trip_id")
	}

	body, err := ParseAndValidateRequestBody[request.TrainJourney](ctx, r.v)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "parse request body")
	}

	transportation, err := r.uc.CreateTrainJourney(ctx.Context(), int32(tripID), *body)
	if err != nil {
		return fmt.Errorf("retrieve journey: %w", err)
	}

	return ctx.Status(http.StatusOK).JSON(transportation)
}
