package app

import (
	"context"
	"errors"
	"fmt"
	"inventory-service/config"
	"inventory-service/internal/adapter/repository"
	"inventory-service/internal/domain/service"
	"inventory-service/pkg/apmtracer"
	"inventory-service/pkg/bundb"
	"inventory-service/pkg/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"inventory-service/proto/pb"

	"google.golang.org/grpc"
)

type App struct {
	config     *config.Config
	grpcServer grpc.Server
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

	var (
		wg  sync.WaitGroup
		err error
	)

	// Initialize tracer
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

	// Initialize repository
	repo, err := repository.NewRepository(a.config, a.logger)
	if err != nil {
		return fmt.Errorf("failed to setup repository: %w", err)
	}

	// Initialize service
	service, err := service.NewService(a.config, repo, a.logger)
	if err != nil {
		return fmt.Errorf("failed to setup service: %w", err)
	}

	// Initialize and start gRPC server with sophisticated setup
	grpcServer := grpc.NewServer()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.config.Grpc.Host, a.config.Grpc.Port))
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			a.logger.Error().Err(err).Msg("Failed to start gRPC server")
			cancel()
		}
	}()

	a.logger.Info().Msgf("Server started at %s:%d", a.config.Grpc.Host, a.config.Grpc.Port)

	// Register gRPC services with proper implementations
	pb.RegisterInventoryServiceServer(grpcServer, service)

	// Wait for shutdown signal
	<-ctx.Done()
	a.logger.Info().Msg("Shutdown signal received, starting graceful shutdown...")

	// Shutdown gRPC server
	if err := grpcServer.Shutdown(); err != nil {
		a.logger.Error().Err(err).Msg("Failed to gracefully shutdown gRPC server")
	} else {
		a.logger.Info().Msg("gRPC server shut down gracefully")
	}

	// Wait for background tasks to finish
	a.logger.Info().Msg("Waiting for background tasks to finish...")
	wg.Wait()
	a.logger.Info().Msg("All background tasks finished")

	// Close repository
	if err := repo.Close(); err != nil {
		a.logger.Error().Err(err).Msg("Failed to gracefully close repository")
	} else {
		a.logger.Info().Msg("Repository closed gracefully")
	}

	a.tracer.Shutdown()

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
