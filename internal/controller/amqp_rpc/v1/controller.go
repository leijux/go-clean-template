package v1

import (
	"github.com/go-playground/validator/v10"
	"github.com/leijux/go-clean-template/internal/usecase"
	"github.com/leijux/go-clean-template/pkg/jwt"
	"github.com/leijux/go-clean-template/pkg/logger"
)

// V1 -.
type V1 struct {
	t  usecase.Translation
	u  usecase.User
	tk usecase.Task
	j  *jwt.Manager
	l  logger.Interface
	v  *validator.Validate
}
