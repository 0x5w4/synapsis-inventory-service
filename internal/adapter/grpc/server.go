package grpcadapter

import (
	"context"
	"fmt"
	"inventory-service/config"
	"inventory-service/internal/adapter/repository"
	"inventory-service/pkg/logger"
	"inventory-service/proto/pb"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

// chainUnaryInterceptors chains multiple gRPC interceptors into one
func chainUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			current := interceptors[i]
			chain = func(currentHandler grpc.UnaryHandler) grpc.UnaryHandler {
				return func(currentCtx context.Context, currentReq any) (any, error) {
					return current(currentCtx, currentReq, info, currentHandler)
				}
			}(chain)
		}
		return chain(ctx, req)
	}
}

func NewGRPCServer(config *config.Config, repo repository.Repository, logger logger.Logger) (*grpc.Server, error) {
	grpcService, err := NewGRPCService(config, repo, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to setup gRPC service: %w", err)
	}

	// Chain logging and tracing interceptors
	interceptor := chainUnaryInterceptors(
		LoggingInterceptor(logger),
		TracingInterceptor(),
	)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor),
	)

	pb.RegisterInventoryServiceServer(grpcServer, grpcService)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Grpc.Host, config.Grpc.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	// Graceful shutdown handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			logger.Error().Err(err).Msg("Failed to start gRPC server")
		}
	}()

	go func() {
		<-stop
		logger.Info().Msg("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	return grpcServer, nil
}
