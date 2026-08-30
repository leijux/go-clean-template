package v1

import (
	"github.com/gofiber/fiber/v3"
	"github.com/leijux/go-clean-template/internal/controller/restapi/v1/response"
)

// okResponse wraps a successful payload in the unified envelope.
func okResponse(ctx fiber.Ctx, code int, data any) error {
	return ctx.Status(code).JSON(response.Envelope{
		Code:    code,
		Message: "ok",
		Data:    data,
	})
}

// okResponseNoData wraps a successful response that carries no payload.
func okResponseNoData(ctx fiber.Ctx, code int) error {
	return okResponse(ctx, code, nil)
}

// errorResponse wraps an error summary in the unified envelope.
func errorResponse(ctx fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(response.Envelope{
		Code:    code,
		Message: msg,
		Data:    nil,
	})
}
