package v1

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"kompass/internal/controller/http/v1/response"
	"kompass/internal/entity"
	netHttp "net/http"
)

func ParseAndValidateRequestBody[V interface{}](ctx *fiber.Ctx, v *validator.Validate) (*V, error) {
	var body V

	if err := ctx.BodyParser(&body); err != nil {
		return nil, fmt.Errorf("parse json body: %w", err)
	}

	if err := v.Struct(body); err != nil {
		return nil, fmt.Errorf("validate json body: %w", err)
	}

	return &body, nil
}

func ErrorResponse(ctx *fiber.Ctx, err error) error {
	return ErrorResponseWithStatus(ctx, netHttp.StatusInternalServerError, err)
}

func ForbiddenResponse(ctx *fiber.Ctx) error {
	return ErrorResponseWithStatus(ctx, netHttp.StatusForbidden, fmt.Errorf("forbidden"))
}

func ErrorResponseWithStatus(ctx *fiber.Ctx, code int, err error) error {
	fmt.Println(err)
	return ctx.Status(code).JSON(response.Error{Error: err.Error()})
}

func userIdFromCtx(ctx *fiber.Ctx) int32 {
	user := ctx.Locals("user").(entity.User)
	return user.ID
}
