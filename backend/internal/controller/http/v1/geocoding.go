package v1

import (
	"fmt"
	"kompass/internal/usecase"
	"kompass/pkg/logger"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type GeocodingV1 struct {
	uc  usecase.Geocoding
	log logger.Interface
	v   *validator.Validate
}

// @Summary     Lookup location
// @ID          getLocation
// @Tags  	    geocoding
// @Produce     json
// @Param       query query string true "location query"
// @Success     200 {object} entity.Location
// @Failure     500 {object} response.Error
// @Router      /geocoding/location [get]
func (r *GeocodingV1) getLocation(ctx *fiber.Ctx) error {
	query := ctx.Query("query")
	location, err := r.uc.LookupLocation(ctx.Context(), query)
	if err != nil {
		return fmt.Errorf("lookup trainstation: %w", err)
	}

	return ctx.Status(http.StatusOK).JSON(location)
}

// @Summary     Lookup train station
// @ID          getTrainStation
// @Tags  	    geocoding
// @Produce     json
// @Param       query query string true "station query"
// @Success     200 {object} entity.TrainStation
// @Failure     500 {object} response.Error
// @Router      /geocoding/station [get]
func (r *GeocodingV1) getTrainStation(ctx *fiber.Ctx) error {
	query := ctx.Query("query")
	location, err := r.uc.LookupTrainStation(ctx.Context(), query)
	if err != nil {
		return fmt.Errorf("lookup trainstation: %w", err)
	}

	return ctx.Status(http.StatusOK).JSON(location)
}
