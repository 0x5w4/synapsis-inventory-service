package grpcadapter

import (
	"context"
	"inventory-service/pkg/logger"

	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"
)

func TracingInterceptor() grpc.UnaryServerInterceptor {
	return apmgrpc.NewUnaryServerInterceptor()
}

func LoggingInterceptor(appLogger logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		appLogger.Info().Field("method", info.FullMethod).Msg("Incoming gRPC request")
		resp, err := handler(ctx, req)
		if err != nil {
			appLogger.Error().Field("method", info.FullMethod).Err(err).Msg("gRPC request failed")
		}
		return resp, err
	}
}
