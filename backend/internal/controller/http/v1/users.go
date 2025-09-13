package v1

import (
	"fmt"
	"kompass/internal/usecase"
	"kompass/pkg/logger"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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
		return errorResponseWithStatus(ctx, http.StatusInternalServerError, fmt.Errorf("internal server error: %w", err))
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
		return errorResponseWithStatus(ctx, http.StatusBadRequest, fmt.Errorf("unable to parse user_id: %w", err))
	}

	user, err := r.uc.GetUserByID(ctx.UserContext(), int32(userID))
	if err != nil {
		return errorResponseWithStatus(ctx, http.StatusNotFound, fmt.Errorf("unable to find user [id=%d]: %w", userID, err))
	}

	return ctx.Status(http.StatusOK).JSON(user)
}
