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
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/activities [get]
func (r *ActivitiesV1) getActivities(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "unable to parse trip_id")
	}

	activities, err := r.uc.GetActivities(ctx.UserContext(), int32(tripID))
	if err != nil {
		r.log.Error(err, "http - v1 - getActivities")
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(activities)
}

// @Summary     Get activity by ID
// @ID          getActivity
// @Tags  	    activities
// @Produce     json
// @Param       trip_id path int true "Trip ID"
// @Param       activity_id path string true "Activity ID"
// @Success     200 {object} entity.Activity
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/activities/{activity_id} [get]
func (r *ActivitiesV1) getActivity(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "unable to parse trip_id")
	}
	activityID, err := ctx.ParamsInt("activity_id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "unable to parse activity_id")
	}

	activity, err := r.uc.GetActivityByID(ctx.UserContext(), int32(tripID), int32(activityID))
	if err != nil {
		return errorResponse(ctx, http.StatusNotFound, fmt.Sprintf("unable to find activity with id %d", activityID))
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
// @Failure     500 {object} response.Error
// @Router      /trips/{trip_id}/activities [post]
func (r *ActivitiesV1) postActivity(ctx *fiber.Ctx) error {
	tripID, err := ctx.ParamsInt("trip_id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "unable to parse trip_id")
	}

	var body request.Activity

	if err := ctx.BodyParser(&body); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		fmt.Println(err)
		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	_, err = r.uc.CreateActivity(ctx.UserContext(), int32(tripID), body)
	if err != nil {
		return err
	}

	return ctx.SendStatus(http.StatusOK)
}
