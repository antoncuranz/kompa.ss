package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	v1 "kompass/internal/controller/http/v1"
	"kompass/internal/entity"
	"kompass/internal/usecase"
)

func TripAuthorization(uc usecase.Users) func(c *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		userID := ctx.Locals("user").(entity.User).ID
		tripID, err := ctx.ParamsInt("trip_id")
		if err != nil {
			return v1.ErrorResponse(ctx, fmt.Errorf("parse trip_id: %w", err))
		}

		switch ctx.Method() {
		case fiber.MethodPost, fiber.MethodPut, fiber.MethodDelete:
			hasPermission, err := uc.HasWritePermission(ctx.UserContext(), userID, int32(tripID))
			if err != nil {
				return v1.ErrorResponse(ctx, fmt.Errorf("has write permission: %w", err))
			}
			if !hasPermission {
				return v1.ForbiddenResponse(ctx)
			}
		default:
			hasPermission, err := uc.HasReadPermission(ctx.UserContext(), userID, int32(tripID))
			if err != nil {
				return v1.ErrorResponse(ctx, fmt.Errorf("has read permission: %w", err))
			}
			if !hasPermission {
				return v1.ForbiddenResponse(ctx)
			}
		}

		return ctx.Next()
	}
}
