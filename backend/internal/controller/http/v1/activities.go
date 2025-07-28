package v1

import (
	"backplate/internal/usecase"
	"backplate/pkg/logger"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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
// @Success     200 {object} []entity.Activity
// @Failure     500 {object} response.Error
// @Router      /activities [get]
func (r *ActivitiesV1) getActivities(ctx *fiber.Ctx) error {
	activities, err := r.uc.GetActivities(ctx.UserContext())
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
// @Param       activity_id path string true "Activity ID"
// @Success     200 {object} entity.Activity
// @Failure     500 {object} response.Error
// @Router      /activities/{activity_id} [get]
func (r *ActivitiesV1) getActivity(ctx *fiber.Ctx) error {
	activityID, err := ctx.ParamsInt("activity_id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "unable to parse activity_id")
	}

	activity, err := r.uc.GetActivityByID(ctx.UserContext(), int32(activityID))
	if err != nil {
		return errorResponse(ctx, http.StatusNotFound, fmt.Sprintf("unable to find activity with id %d", activityID))
	}

	return ctx.Status(http.StatusOK).JSON(activity)
}
