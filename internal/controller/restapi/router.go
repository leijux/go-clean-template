package restapi

import (
	"log/slog"
	"net/http"

	"github.com/gofiber/contrib/v3/otel"
	"github.com/gofiber/contrib/v3/prometheus"
	"github.com/gofiber/contrib/v3/swaggerui"
	"github.com/gofiber/fiber/v3"
	"github.com/leijux/go-clean-template/config"
	_ "github.com/leijux/go-clean-template/docs" // Swagger docs.
	"github.com/leijux/go-clean-template/internal/controller/restapi/middleware"
	v1 "github.com/leijux/go-clean-template/internal/controller/restapi/v1"
	"github.com/leijux/go-clean-template/internal/usecase"
)

// NewRouter -.
// Swagger spec:
//
//	@title       Go Clean Template API
//	@description Multi-domain clean architecture template with translation, user, and task management
//	@version     1.0
//	@host        localhost:8080
//	@BasePath    /v1
//	@securityDefinitions.apikey BearerAuth
//	@in header
//	@name Authorization
func NewRouter(app *fiber.App, cfg *config.Config, t usecase.Translation, u usecase.User, tk usecase.Task, l *slog.Logger) {
	// Options
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))

	// Prometheus metrics
	if cfg.Metrics.Enabled {
		app.Use(prometheus.New(prometheus.Config{ServiceName: cfg.App.Name}))
	}

	// Swagger
	if cfg.Swagger.Enabled {
		app.Use(swaggerui.New(swaggerui.Config{Path: "swagger", FilePath: "./docs/swagger.json"}))
	}

	// K8s probe
	app.Get("/healthz", func(ctx fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	// Routers
	apiV1Group := app.Group("/v1")
	{
		if cfg.Tracing.Enabled {
			apiV1Group.Use(otel.Middleware())
		}

		v1.NewRoutes(apiV1Group, t, u, tk, []byte(cfg.JWT.Secret), l)
	}
}
