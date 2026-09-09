// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/leijux/go-clean-template/config"
	amqprpc "github.com/leijux/go-clean-template/internal/controller/amqp_rpc"
	"github.com/leijux/go-clean-template/internal/controller/grpc"
	grpcmw "github.com/leijux/go-clean-template/internal/controller/grpc/middleware"
	natsrpc "github.com/leijux/go-clean-template/internal/controller/nats_rpc"
	"github.com/leijux/go-clean-template/internal/controller/restapi"
	cacheTranslationRepo "github.com/leijux/go-clean-template/internal/repo/cache/translation"
	persistTaskRepo "github.com/leijux/go-clean-template/internal/repo/persistent/task"
	persistTranslationRepo "github.com/leijux/go-clean-template/internal/repo/persistent/translation"
	persistUserRepo "github.com/leijux/go-clean-template/internal/repo/persistent/user"
	"github.com/leijux/go-clean-template/internal/repo/webapi"
	"github.com/leijux/go-clean-template/internal/usecase"
	"github.com/leijux/go-clean-template/internal/usecase/task"
	"github.com/leijux/go-clean-template/internal/usecase/translation"
	"github.com/leijux/go-clean-template/internal/usecase/user"
	"github.com/leijux/go-clean-template/pkg/grpcserver"
	"github.com/leijux/go-clean-template/pkg/httpserver"
	"github.com/leijux/go-clean-template/pkg/jwt"
	"github.com/leijux/go-clean-template/pkg/logger"
	natsRPCServer "github.com/leijux/go-clean-template/pkg/nats/nats_rpc/server"
	"github.com/leijux/go-clean-template/pkg/postgres"
	rmqRPCServer "github.com/leijux/go-clean-template/pkg/rabbitmq/rmq_rpc/server"
	"github.com/leijux/go-clean-template/pkg/redis"
	"github.com/leijux/go-clean-template/pkg/tracing"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	pbgrpc "google.golang.org/grpc"
)

type useCases struct {
	translation usecase.Translation
	user        usecase.User
	task        usecase.Task
}

type servers struct {
	rmq  *rmqRPCServer.Server
	nats *natsRPCServer.Server
	grpc *grpcserver.Server
	http *httpserver.Server
}

func initUseCases(pg *postgres.Postgres, rdb *redis.Redis, jwtManager *jwt.Manager) useCases {
	translationRepo := persistTranslationRepo.New(pg)
	translationCache := cacheTranslationRepo.New(rdb.Client)
	taskRepo := persistTaskRepo.New(pg)
	userRepo := persistUserRepo.New(pg)

	return useCases{
		user:        user.New(userRepo, jwtManager),
		task:        task.New(taskRepo),
		translation: translation.New(translationRepo, webapi.New(), translationCache),
	}
}

func initServers(cfg *config.Config, uc useCases, jwtManager *jwt.Manager, l *slog.Logger) (servers, error) {
	// RabbitMQ RPC Server
	rmqRouter := amqprpc.NewRouter(uc.translation, uc.user, uc.task, jwtManager, l)

	rmqServer, err := rmqRPCServer.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, l)
	if err != nil {
		return servers{}, fmt.Errorf("app - Run - rmqServer - server.New: %w", err)
	}

	// NATS RPC Server
	natsRouter := natsrpc.NewRouter(uc.translation, uc.user, uc.task, jwtManager, l)

	natsServer, err := natsRPCServer.New(cfg.NATS.URL, cfg.NATS.ServerExchange, natsRouter, l)
	if err != nil {
		return servers{}, fmt.Errorf("app - Run - natsServer - server.New: %w", err)
	}

	// gRPC Server
	grpcServer := grpcserver.New(
		l,
		grpcserver.Port(cfg.GRPC.Port),
		grpcserver.ServerOptions(
			pbgrpc.UnaryInterceptor(grpcmw.AuthInterceptor(jwtManager)),
			pbgrpc.StatsHandler(otelgrpc.NewServerHandler()),
		),
	)
	grpc.NewRouter(grpcServer.App, uc.translation, uc.user, uc.task, l)

	// HTTP Server
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	restapi.NewRouter(httpServer.App, cfg, uc.translation, uc.user, uc.task, l)

	return servers{
		rmq:  rmqServer,
		nats: natsServer,
		grpc: grpcServer,
		http: httpServer,
	}, nil
}

func (s *servers) startServers() {
	s.rmq.Start()
	s.nats.Start()
	s.grpc.Start()
	s.http.Start()
}

func (s *servers) waitForShutdown(l *slog.Logger) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error

	select {
	case sig := <-interrupt:
		l.Info("app - Run - signal", "signal", sig.String())
	case err = <-s.http.Notify():
		l.Error("app - Run - httpServer.Notify", "error", err)
	case err = <-s.grpc.Notify():
		l.Error("app - Run - grpcServer.Notify", "error", err)
	case err = <-s.rmq.Notify():
		l.Error("app - Run - rmqServer.Notify", "error", err)
	case err = <-s.nats.Notify():
		l.Error("app - Run - natsServer.Notify", "error", err)
	}

	s.shutdownServers(l)
}

func (s *servers) shutdownServers(l *slog.Logger) {
	if err := s.http.Shutdown(); err != nil {
		l.Error("app - Run - httpServer.Shutdown", "error", err)
	}

	if err := s.grpc.Shutdown(); err != nil {
		l.Error("app - Run - grpcServer.Shutdown", "error", err)
	}

	if err := s.rmq.Shutdown(); err != nil {
		l.Error("app - Run - rmqServer.Shutdown", "error", err)
	}

	if err := s.nats.Shutdown(); err != nil {
		l.Error("app - Run - natsServer.Shutdown", "error", err)
	}
}

// Run creates objects via constructors.
func Run(cfg *config.Config) error {
	l := logger.New(cfg.Log.Level, cfg.Log.File)

	ctx := context.Background()

	// Tracing
	shutdownTracing, err := tracing.New(ctx, tracing.Config{
		Enabled:     cfg.Tracing.Enabled,
		ServiceName: cfg.App.Name,
		Version:     cfg.App.Version,
		Endpoint:    cfg.Tracing.OTLPEndpoint,
		Insecure:    cfg.Tracing.OTLPInsecure,
		SampleRate:  cfg.Tracing.SampleRate,
	})
	if err != nil {
		return fmt.Errorf("app - Run - tracing.New: %w", err)
	}
	defer func() {
		if err := shutdownTracing(ctx); err != nil {
			l.Error("app - Run - shutdownTracing", "error", err)
		}
	}()

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		return fmt.Errorf("app - Run - postgres.New: %w", err)
	}
	defer pg.Close()

	// Cache
	rdb, err := redis.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		return fmt.Errorf("app - Run - redis.New: %w", err)
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			l.Error("app - Run - redis.Close", "error", err)
		}
	}()

	// JWT
	jwtManager := jwt.New(cfg.JWT.Secret, cfg.JWT.TokenExpiry)

	uc := initUseCases(pg, rdb, jwtManager)

	s, err := initServers(cfg, uc, jwtManager, l)
	if err != nil {
		return err
	}

	s.startServers()
	s.waitForShutdown(l)

	return nil
}
