package v1

import (
	"log/slog"

	v1 "github.com/leijux/go-clean-template/internal/controller/amqp_rpc/v1"
	"github.com/leijux/go-clean-template/internal/usecase"
	"github.com/leijux/go-clean-template/pkg/jwt"
	"github.com/leijux/go-clean-template/pkg/rabbitmq/rmq_rpc/server"
)

// NewRouter -.
func NewRouter(t usecase.Translation, u usecase.User, tk usecase.Task, j *jwt.Manager, l *slog.Logger) map[string]server.CallHandler {
	routes := make(map[string]server.CallHandler)

	{
		v1.NewRoutes(routes, t, u, tk, j, l)
	}

	return routes
}
