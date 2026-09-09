package v1

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	v1 "github.com/leijux/go-clean-template/docs/proto/v1"
	"github.com/leijux/go-clean-template/internal/usecase"
)

// AuthController -.
type AuthController struct {
	v1.UnimplementedAuthServiceServer

	u usecase.User
	l *slog.Logger
	v *validator.Validate
}
