package app

import (
	"context"
	"errors"
	"fmt"
	"inventory-service/config"
	grpcadapter "inventory-service/internal/adapter/grpc"
	"inventory-service/internal/adapter/repository"
	"inventory-service/pkg/apmtracer"
	"inventory-service/pkg/bundb"
	"inventory-service/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

type App struct {
	config     *config.Config
	grpcServer *grpc.Server
	logger     logger.Logger
	tracer     apmtracer.Tracer
}

func NewApp(config *config.Config, logger logger.Logger) (*App, error) {
	if config == nil {
		return nil, errors.New("configuration cannot be nil")
	}

	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}

	return &App{
		config: config,
		logger: logger,
	}, nil
}

func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	var err error

	a.tracer, err = apmtracer.NewApmTracer(&apmtracer.Config{
		ServiceName:    a.config.Tracer.ServiceName,
		ServiceVersion: a.config.Tracer.ServiceVersion,
		ServerURL:      a.config.Tracer.ServerURL,
		SecretToken:    a.config.Tracer.SecretToken,
		Environment:    a.config.Tracer.Environment,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize tracer: %w", err)
	}

	repo, err := repository.NewRepository(a.config, a.logger)
	if err != nil {
		return fmt.Errorf("failed to setup repository: %w", err)
	}
	defer func() {
		if err := repo.Close(); err != nil {
			a.logger.Error().Err(err).Msg("Failed to gracefully close repository")
		} else {
			a.logger.Info().Msg("Repository closed gracefully")
		}
	}()

	a.grpcServer, err = grpcadapter.NewGRPCServer(a.config, repo, a.logger)
	if err != nil {
		return fmt.Errorf("failed to setup gRPC server: %w", err)
	}

	a.logger.Info().Msgf("Server started at %s:%d", a.config.Grpc.Host, a.config.Grpc.Port)

	<-ctx.Done()
	a.logger.Info().Msg("Shutdown signal received, starting graceful shutdown...")
	a.grpcServer.GracefulStop()
	a.logger.Info().Msg("gRPC server shut down gracefully")

	if a.tracer != nil {
		a.tracer.Shutdown()
	}

	return nil
}

func (a *App) Migrate(reset bool) error {
	db, err := bundb.NewBunDB(a.config, a.logger)
	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			a.logger.Error().Err(err).Msg("Failed to close database connection")
		}
	}()

	if reset {
		if err := db.Reset(); err != nil {
			return err
		}
	} else {
		if err := db.Migrate(); err != nil {
			return err
		}
	}

	return nil
}
