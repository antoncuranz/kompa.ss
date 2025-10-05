package v1

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"kompass/internal/entity"
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

func userIdFromCtx(ctx *fiber.Ctx) int32 {
	user := ctx.Locals("user").(entity.User)
	return user.ID
}
