package v1

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	v1 "github.com/leijux/go-clean-template/docs/proto/v1"
	"github.com/leijux/go-clean-template/internal/usecase"
)

// TaskController -.
type TaskController struct {
	v1.UnimplementedTaskServiceServer

	tk usecase.Task
	l  *slog.Logger
	v  *validator.Validate
}
