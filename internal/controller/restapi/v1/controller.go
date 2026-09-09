package v1

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/leijux/go-clean-template/internal/usecase"
)

// V1 -.
type V1 struct {
	t  usecase.Translation
	u  usecase.User
	tk usecase.Task
	l  *slog.Logger
	v  *validator.Validate
}
