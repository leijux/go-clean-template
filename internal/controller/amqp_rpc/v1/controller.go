package v1

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/leijux/go-clean-template/internal/usecase"
	"github.com/leijux/go-clean-template/pkg/jwt"
)

// V1 -.
type V1 struct {
	t  usecase.Translation
	u  usecase.User
	tk usecase.Task
	j  *jwt.Manager
	l  *slog.Logger
	v  *validator.Validate
}
