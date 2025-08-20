package v1

import (
	"kompass/pkg/logger"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TransportationV1 struct {
	//uc  usecase.Transportation
	log logger.Interface
	v   *validator.Validate
}

// @Summary     Get all Transportation
// @ID          getAllTransportation
// @Tags  	    transportation
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Success     200 {object} []response.Transportation
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/Transportation [get]
func (r *TransportationV1) getAllTransportation(ctx *fiber.Ctx) error {
	return ctx.SendStatus(http.StatusOK)
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
	return ctx.SendStatus(http.StatusOK) //.JSON(Transportation)
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
	return ctx.SendStatus(http.StatusNoContent)
}
