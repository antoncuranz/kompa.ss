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

type TransportationV1 struct {
	uc  usecase.Transportation
	log logger.Interface
	v   *validator.Validate
}

// @Summary     Add transportation
// @ID          postTransportation
// @Tags  	    transportation
// @Accept      json
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       request body request.Transportation true "transportation"
// @Success     204
// @Failure     403 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id}/transportation [post]
func (r *TransportationV1) postTransportation(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return ErrorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse trip_id: %w", err))
	}

	body, err := ParseAndValidateRequestBody[request.Transportation](ctx, r.v)
	if err != nil {
		return ErrorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
	}

	_, err = r.uc.CreateTransportation(ctx.UserContext(), int32(tripID), *body)
	if err != nil {
		return ErrorResponse(ctx, fmt.Errorf("create transportation: %w", err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     Update transportation
// @ID          putTransportation
// @Tags  	    transportation
// @Accept      json
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       transportation_id path int true "Transportation ID"
// @Param       request body request.Transportation true "transportation"
// @Success     204
// @Failure     403 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id}/transportation/{transportation_id} [put]
func (r *TransportationV1) putTransportation(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return ErrorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse trip_id: %w", err))
	}
	transportationID, err := ctx.ParamsInt("transportation_id")
	if err != nil {
		return ErrorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse transportation_id: %w", err))
	}

	body, err := ParseAndValidateRequestBody[request.Transportation](ctx, r.v)
	if err != nil {
		return ErrorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
	}

	if _, err := r.uc.UpdateTransportation(ctx.UserContext(), int32(tripID), int32(transportationID), *body); err != nil {
		return ErrorResponse(ctx, fmt.Errorf("update transportation with id %d: %w", transportationID, err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     Get all Transportation
// @ID          getAllTransportation
// @Tags  	    transportation
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Success     200 {object} []entity.Transportation
// @Failure     403 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id}/transportation [get]
func (r *TransportationV1) getAllTransportation(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return ErrorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("parse trip_id: %w", err))
	}

	transportation, err := r.uc.GetAllTransportation(ctx.Context(), int32(tripID))
	if err != nil {
		return ErrorResponse(ctx, fmt.Errorf("get all transportation: %w", err))
	}

	return ctx.Status(http.StatusOK).JSON(transportation)
}

// @Summary     Get Transportation by ID
// @ID          getTransportation
// @Tags  	    transportation
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       transportation_id path int true "Transportation ID"
// @Success     200 {object} entity.Transportation
// @Failure     403 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id}/transportation/{transportation_id} [get]
func (r *TransportationV1) getTransportation(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return ErrorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("parse trip_id: %w", err))
	}
	transportationID, err := ctx.ParamsInt("transportation_id")
	if err != nil {
		return ErrorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("parse transportation_id: %w", err))
	}

	transportation, err := r.uc.GetTransportationByID(ctx.UserContext(), int32(tripID), int32(transportationID))
	if err != nil {
		return ErrorResponse(ctx, fmt.Errorf("get transportation [id=%d]: %w", transportationID, err))
	}

	return ctx.Status(http.StatusOK).JSON(transportation)
}

// @Summary     Delete Transportation
// @ID          deleteTransportation
// @Tags  	    transportation
// @Param       trip_id path int true "Trip ID"
// @Param       transportation_id path int true "Transportation ID"
// @Success     204
// @Failure     403 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id}/transportation/{transportation_id} [delete]
func (r *TransportationV1) deleteTransportation(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return ErrorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("parse trip_id: %w", err))
	}
	transportationID, err := ctx.ParamsInt("transportation_id")
	if err != nil {
		return ErrorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("parse transportation_id: %w", err))
	}

	if err := r.uc.DeleteTransportation(ctx.UserContext(), int32(tripID), int32(transportationID)); err != nil {
		return ErrorResponse(ctx, fmt.Errorf("delete transportation with id %d: %w", transportationID, err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     Get GeoJson
// @ID          getGeoJson
// @Tags  	    transportation
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Success     200 {object} []geojson.FeatureCollection
// @Failure     403 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id}/transportation/geojson [get]
func (r *TransportationV1) getGeoJson(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return ErrorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("parse trip_id: %w", err))
	}

	geojson, err := r.uc.GetAllGeoJson(ctx.Context(), int32(tripID))
	if err != nil {
		return ErrorResponse(ctx, fmt.Errorf("get geojson: %w", err))
	}

	return ctx.Status(http.StatusOK).JSON(geojson)
}
