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

type AccommodationV1 struct {
	uc  usecase.Accommodation
	log logger.Interface
	v   *validator.Validate
}

// @Summary     Get all accommodation
// @ID          getAllAccommodation
// @Tags  	    accommodation
// @Produce     json
// @Success     200 {object} []entity.Accommodation
// @Failure     500 {object} response.Error
// @Router      /accommodation [get]
func (r *AccommodationV1) getAllAccommodation(ctx *fiber.Ctx) error {
	accommodation, err := r.uc.GetAllAccommodation(ctx.UserContext())
	if err != nil {
		r.log.Error(err, "http - v1 - getAccommodation")
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(accommodation)
}

// @Summary     Get accommodation by ID
// @ID          getAccommodationByID
// @Tags  	    accommodation
// @Produce     json
// @Param       accommodation_id path string true "Accommodation ID"
// @Success     200 {object} entity.Accommodation
// @Failure     500 {object} response.Error
// @Router      /accommodation/{accommodation_id} [get]
func (r *AccommodationV1) getAccommodationByID(ctx *fiber.Ctx) error {
	accommodationID, err := ctx.ParamsInt("accommodation_id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "unable to parse accommodation_id")
	}

	accommodation, err := r.uc.GetAccommodationByID(ctx.UserContext(), int32(accommodationID))
	if err != nil {
		return errorResponse(ctx, http.StatusNotFound, fmt.Sprintf("unable to find accommodation with id %d", accommodationID))
	}

	return ctx.Status(http.StatusOK).JSON(accommodation)
}

// @Summary     Add accommodation
// @ID          postAccommodation
// @Tags  	    accommodation
// @Accept      json
// @Produce     json
// @Param       request body request.Accommodation true "accommodation"
// @Success     200 {object} entity.Accommodation
// @Failure     500 {object} response.Error
// @Router      /accommodation [post]
func (r *AccommodationV1) postAccommodation(ctx *fiber.Ctx) error {
	var body request.Accommodation

	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		fmt.Println(err)
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	_, err := r.uc.CreateAccommodation(ctx.UserContext(), body)
	if err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusOK)
}
