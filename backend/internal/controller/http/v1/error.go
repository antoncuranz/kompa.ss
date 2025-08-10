package v1

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"travel-planner/internal/controller/http/v1/response"
)

func errorResponseDeprecated(ctx *fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(response.Error{Error: msg})
}

func errorResponse(ctx *fiber.Ctx, err error) error {
	return errorResponseWithStatus(ctx, http.StatusInternalServerError, err)
}

func errorResponseWithStatus(ctx *fiber.Ctx, code int, err error) error {
	fmt.Println(err)
	return ctx.Status(code).JSON(response.Error{Error: err.Error()})
}
