package v1

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	v1 "github.com/leijux/go-clean-template/docs/proto/v1"
	"github.com/leijux/go-clean-template/internal/usecase"
)

// TranslationController -.
type TranslationController struct {
	v1.UnimplementedTranslationServer

	t usecase.Translation
	l *slog.Logger
	v *validator.Validate
}
