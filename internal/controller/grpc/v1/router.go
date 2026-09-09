package v1

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	v1 "github.com/leijux/go-clean-template/docs/proto/v1"
	"github.com/leijux/go-clean-template/internal/usecase"
	pbgrpc "google.golang.org/grpc"
)

// NewTranslationRoutes -.
func NewTranslationRoutes(app *pbgrpc.Server, t usecase.Translation, l *slog.Logger) {
	r := &TranslationController{t: t, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	v1.RegisterTranslationServer(app, r)
}

// NewAuthRoutes -.
func NewAuthRoutes(app *pbgrpc.Server, u usecase.User, l *slog.Logger) {
	r := &AuthController{u: u, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	v1.RegisterAuthServiceServer(app, r)
}

// NewTaskRoutes -.
func NewTaskRoutes(app *pbgrpc.Server, tk usecase.Task, l *slog.Logger) {
	r := &TaskController{tk: tk, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

	v1.RegisterTaskServiceServer(app, r)
}
