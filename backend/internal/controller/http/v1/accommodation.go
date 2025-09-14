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

type AccommodationV1 struct {
	uc  usecase.Accommodation
	log logger.Interface
	v   *validator.Validate
}

// @Summary     Get all accommodation
// @ID          getAllAccommodation
// @Tags  	    accommodation
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Success     200 {object} []entity.Accommodation
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/accommodation [get]
func (r *AccommodationV1) getAllAccommodation(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse trip_id: %v", err))
	}

	accommodation, err := r.uc.GetAllAccommodation(ctx.UserContext(), int32(tripID))
	if err != nil {
		return errorResponse(ctx, fmt.Errorf("get accommodation: %w", err))
	}

	return ctx.Status(http.StatusOK).JSON(accommodation)
}

// @Summary     Get accommodation by ID
// @ID          getAccommodationByID
// @Tags  	    accommodation
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       accommodation_id path int true "Accommodation ID"
// @Success     200 {object} entity.Accommodation
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/accommodation/{accommodation_id} [get]
func (r *AccommodationV1) getAccommodationByID(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse trip_id: %v", err))
	}
	accommodationID, err := ctx.ParamsInt("accommodation_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse accommodation_id: %w", err))
	}

	accommodation, err := r.uc.GetAccommodationByID(ctx.UserContext(), int32(tripID), int32(accommodationID))
	if err != nil {
		return errorResponse(ctx, fmt.Errorf("get accommodation [id=%d]: %w", accommodationID, err))
	}

	return ctx.Status(http.StatusOK).JSON(accommodation)
}

// @Summary     Add accommodation
// @ID          postAccommodation
// @Tags  	    accommodation
// @Accept      json
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       request body request.Accommodation true "accommodation"
// @Success     204
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/accommodation [post]
func (r *AccommodationV1) postAccommodation(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse trip_id: %w", err))
	}

	body, err := ParseAndValidateRequestBody[request.Accommodation](ctx, r.v)
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
	}

	_, err = r.uc.CreateAccommodation(ctx.UserContext(), int32(tripID), *body)
	if err != nil {
		return errorResponse(ctx, fmt.Errorf("create accommodation: %w", err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     Update accommodation
// @ID          putAccommodation
// @Tags  	    accommodation
// @Accept      json
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       accommodation_id path int true "Accommodation ID"
// @Param       request body request.Accommodation true "accommodation"
// @Success     204
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/accommodation/{accommodation_id} [put]
func (r *AccommodationV1) putAccommodation(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse trip_id: %w", err))
	}
	accommodationID, err := ctx.ParamsInt("accommodation_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse accommodation_id: %w", err))
	}

	body, err := ParseAndValidateRequestBody[request.Accommodation](ctx, r.v)
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("invalid request body: %w", err))
	}

	if err := r.uc.UpdateAccommodation(ctx.UserContext(), int32(tripID), int32(accommodationID), *body); err != nil {
		return errorResponse(ctx, fmt.Errorf("update accommodation with id %d: %w", accommodationID, err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     Delete accommodation
// @ID          deleteAccommodation
// @Tags  	    accommodation
// @Param       trip_id path int true "Trip ID"
// @Param       accommodation_id path int true "Accommodation ID"
// @Success     204
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/accommodation/{accommodation_id} [delete]
func (r *AccommodationV1) deleteAccommodation(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse trip_id: %w", err))
	}
	accommodationID, err := ctx.ParamsInt("accommodation_id")
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse accommodation_id: %w", err))
	}

	if err := r.uc.DeleteAccommodation(ctx.UserContext(), int32(tripID), int32(accommodationID)); err != nil {
		return errorResponse(ctx, fmt.Errorf("delete accommodation with id %d: %w", accommodationID, err))
	}

	return ctx.SendStatus(http.StatusNoContent)
}
