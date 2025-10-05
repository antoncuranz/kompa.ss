package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"kompass/internal/entity"
	"kompass/internal/usecase"
	"net/http"
)

func TripAuthorization(uc usecase.Users) func(c *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		userID := ctx.Locals("user").(entity.User).ID
		tripID, err := ctx.ParamsInt("trip_id")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "unable to parse trip_id")
		}

		switch ctx.Method() {
		case fiber.MethodPost, fiber.MethodPut, fiber.MethodDelete:
			hasPermission, err := uc.HasWritePermission(ctx.UserContext(), userID, int32(tripID))
			if err != nil {
				return fmt.Errorf("has write permission: %w", err)
			}
			if !hasPermission {
				return fiber.NewError(http.StatusForbidden)
			}
		default:
			hasPermission, err := uc.HasReadPermission(ctx.UserContext(), userID, int32(tripID))
			if err != nil {
				return fmt.Errorf("has read permission: %w", err)
			}
			if !hasPermission {
				return fiber.NewError(http.StatusForbidden)
			}
		}

		return ctx.Next()
	}
}
