package v1

import (
	"github.com/gofiber/fiber/v3"
	"github.com/leijux/go-clean-template/internal/controller/restapi/v1/response"
)

// okResponse wraps a successful payload in the success envelope.
func okResponse[T any](ctx fiber.Ctx, code int, data T) error {
	return ctx.Status(code).JSON(response.Ok[T]{
		Code:    code,
		Message: "ok",
		Data:    data,
	})
}

// okResponseNoData wraps a successful response that carries no payload.
func okResponseNoData(ctx fiber.Ctx, code int) error {
	return okResponse[any](ctx, code, nil)
}

// errorResponse wraps an error summary in the error envelope.
func errorResponse(ctx fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(response.Error{
		Code:    code,
		Message: msg,
	})
}
