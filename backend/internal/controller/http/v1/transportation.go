package v1

import (
	"fmt"
	"kompass/internal/controller/http/v1/converter"
	"kompass/internal/usecase"
	"kompass/pkg/logger"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TransportationV1 struct {
	uc  usecase.Transportation
	log logger.Interface
	v   *validator.Validate
	c   converter.TransportationConverter
}

// @Summary     Get all Transportation
// @ID          getAllTransportation
// @Tags  	    transportation
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Success     200 {object} []response.Transportation
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/transportation [get]
func (r *TransportationV1) getAllTransportation(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("parse trip_id: %w", err))
	}

	transportation, err := r.uc.GetAllTransportation(ctx.Context(), int32(tripID))
	if err != nil {
		return errorResponse(ctx, fmt.Errorf("get all transportation: %w", err))
	}

	return ctx.Status(http.StatusOK).JSON(r.c.ConvertAllTransportation(transportation))
}

// @Summary     Get Transportation by ID
// @ID          getTransportation
// @Tags  	    transportation
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       transportation_id path int true "Transportation ID"
// @Success     200 {object} response.Transportation
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/transportation/{transportation_id} [get]
func (r *TransportationV1) getTransportation(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("parse trip_id: %w", err))
	}
	transportationID, err := ctx.ParamsInt("transportation_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("parse transportation_id: %w", err))
	}

	transportation, err := r.uc.GetTransportationByID(ctx.UserContext(), int32(tripID), int32(transportationID))
	if err != nil {
		return errorResponse(ctx, fmt.Errorf("get transportation [id=%d]: %w", transportationID, err))
	}

	return ctx.Status(http.StatusOK).JSON(r.c.ConvertTransportation(transportation))
}

// @Summary     Delete Transportation
// @ID          deleteTransportation
// @Tags  	    transportation
// @Param       trip_id path int true "Trip ID"
// @Param       transportation_id path int true "Transportation ID"
// @Success     204
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/transportation/{transportation_id} [delete]
func (r *TransportationV1) deleteTransportation(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("parse trip_id: %w", err))
	}
	transportationID, err := ctx.ParamsInt("transportation_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("parse transportation_id: %w", err))
	}

	if err := r.uc.DeleteTransportation(ctx.UserContext(), int32(tripID), int32(transportationID)); err != nil {
		return errorResponse(ctx, fmt.Errorf("delete transportation with id %d: %w", transportationID, err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}
