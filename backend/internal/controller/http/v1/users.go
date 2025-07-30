package v1

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"travel-planner/internal/usecase"
	"travel-planner/pkg/logger"
)

type UsersV1 struct {
	uc  usecase.Users
	log logger.Interface
	v   *validator.Validate
}

// @Summary     Get all users
// @ID          getUsers
// @Tags  	    users
// @Produce     json
// @Success     200 {object} []entity.User
// @Failure     500 {object} response.Error
// @Router      /users [get]
func (r *UsersV1) getUsers(ctx *fiber.Ctx) error {
	users, err := r.uc.GetUsers(ctx.UserContext())
	if err != nil {
		r.log.Error(err, "http - v1 - getUsers")
		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(users)
}

// @Summary     Get user by ID
// @ID          getUser
// @Tags  	    users
// @Produce     json
// @Param       user_id path int true "User ID"
// @Success     200 {object} entity.User
// @Failure     500 {object} response.Error
// @Router      /users/{user_id} [get]
func (r *UsersV1) getUser(ctx *fiber.Ctx) error {
	userID, err := ctx.ParamsInt("user_id")
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, "unable to parse user_id")
	}

	user, err := r.uc.GetUserByID(ctx.UserContext(), int32(userID))
	if err != nil {
		return errorResponse(ctx, http.StatusNotFound, fmt.Sprintf("unable to find user with id %d", userID))
	}

	return ctx.Status(http.StatusOK).JSON(user)
}
