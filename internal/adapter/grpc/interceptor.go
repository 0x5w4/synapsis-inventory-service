package grpcadapter

import (
	"context"
	"inventory-service/pkg/apmtracer"
	"inventory-service/pkg/logger"

	"go.elastic.co/apm/v2"
	"google.golang.org/grpc"
)

func UnaryInterceptor(appLogger logger.Logger, appTracer apmtracer.Tracer) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		tx := appTracer.Tracer().StartTransaction(info.FullMethod, "request")
		defer tx.End()

		ctx = apm.ContextWithTransaction(ctx, tx)

		appLogger.Info().Field("method", info.FullMethod).Msg("Handling gRPC method")

		resp, err := handler(ctx, req)
		if err != nil {
			appLogger.Error().Field("method", info.FullMethod).Err(err).Msg("gRPC method failed")
		}

		return resp, err
	}
}
