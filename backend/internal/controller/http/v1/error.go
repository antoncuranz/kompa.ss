package v1

import (
	"fmt"
	"kompass/internal/controller/http/v1/response"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func errorResponse(ctx *fiber.Ctx, err error) error {
	return errorResponseWithStatus(ctx, http.StatusInternalServerError, err)
}

func errorResponseWithStatus(ctx *fiber.Ctx, code int, err error) error {
	fmt.Println(err)
	return ctx.Status(code).JSON(response.Error{Error: err.Error()})
}
