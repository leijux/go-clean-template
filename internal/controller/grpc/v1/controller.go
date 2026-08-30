package v1

import (
	"github.com/go-playground/validator/v10"
	v1 "github.com/leijux/go-clean-template/docs/proto/v1"
	"github.com/leijux/go-clean-template/internal/usecase"
	"github.com/leijux/go-clean-template/pkg/logger"
)

// TranslationController -.
type TranslationController struct {
	v1.UnimplementedTranslationServer

	t usecase.Translation
	l logger.Interface
	v *validator.Validate
}
