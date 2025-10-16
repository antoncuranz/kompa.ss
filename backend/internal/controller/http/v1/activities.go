package v1

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"kompass/internal/controller/http/v1/request"
	"kompass/internal/usecase"
	"kompass/pkg/logger"
	"net/http"
)

type ActivitiesV1 struct {
	uc  usecase.Activities
	log logger.Interface
	v   *validator.Validate
}

// @Summary     Get all activities
// @ID          getActivities
// @Tags  	    activities
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Success     200 {object} []entity.Activity
// @Failure     400 {object} response.Error
// @Failure     403 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id}/activities [get]
func (r *ActivitiesV1) getActivities(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "unable to parse trip_id")
	}

	activities, err := r.uc.GetActivities(ctx.UserContext(), int32(tripID))
	if err != nil {
		return fmt.Errorf("get activities: %w", err)
	}

	return ctx.Status(http.StatusOK).JSON(activities)
}

// @Summary     Get activity by ID
// @ID          getActivity
// @Tags  	    activities
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       activity_id path int true "Activity ID"
// @Success     200 {object} entity.Activity
// @Failure     400 {object} response.Error
// @Failure     403 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id}/activities/{activity_id} [get]
func (r *ActivitiesV1) getActivity(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "unable to parse trip_id")
	}
	activityID, err := ctx.ParamsInt("activity_id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "unable to parse activity_id")
	}

	activity, err := r.uc.GetActivityByID(ctx.UserContext(), int32(tripID), int32(activityID))
	if err != nil {
		return fmt.Errorf("get activity [id=%d]: %w", activityID, err)
	}

	return ctx.Status(http.StatusOK).JSON(activity)
}

// @Summary     Add activity
// @ID          postActivity
// @Tags  	    activities
// @Accept      json
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       request body request.Activity true "activity"
// @Success     200 {object} entity.Activity
// @Failure     400 {object} response.Error
// @Failure     403 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id}/activities [post]
func (r *ActivitiesV1) postActivity(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "unable to parse trip_id")
	}

	body, err := ParseAndValidateRequestBody[request.Activity](ctx, r.v)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid request body")
	}

	activity, err := r.uc.CreateActivity(ctx.UserContext(), int32(tripID), *body)
	if err != nil {
		return fmt.Errorf("create activity: %w", err)
	}

	return ctx.Status(http.StatusOK).JSON(activity)
}

// @Summary     Update activity
// @ID          putActivity
// @Tags  	    activities
// @Accept      json
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       activity_id path int true "Activity ID"
// @Param       request body request.Activity true "activity"
// @Success     204
// @Failure     400 {object} response.Error
// @Failure     403 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id}/activities/{activity_id} [put]
func (r *ActivitiesV1) putActivity(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "unable to parse trip_id")
	}
	activityID, err := ctx.ParamsInt("activity_id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "unable to parse activity_id")
	}

	body, err := ParseAndValidateRequestBody[request.Activity](ctx, r.v)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid request body")
	}

	if err := r.uc.UpdateActivity(ctx.UserContext(), int32(tripID), int32(activityID), *body); err != nil {
		return fmt.Errorf("update activity with id %d: %w", activityID, err)
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     Delete activity
// @ID          deleteActivity
// @Tags  	    activities
// @Param       trip_id path int true "Trip ID"
// @Param       activity_id path int true "Activity ID"
// @Success     204
// @Failure     400 {object} response.Error
// @Failure     403 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    bearerauth
// @Router      /trips/{trip_id}/activities/{activity_id} [delete]
func (r *ActivitiesV1) deleteActivity(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "unable to parse trip_id")
	}
	activityID, err := ctx.ParamsInt("activity_id")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "unable to parse activity_id")
	}

	if err := r.uc.DeleteActivity(ctx.UserContext(), int32(tripID), int32(activityID)); err != nil {
		return fmt.Errorf("delete activity with id %d: %w", activityID, err)
	}

	return ctx.SendStatus(http.StatusNoContent)
}
