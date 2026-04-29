package middleware

import (
	"context"
	"time"

	"dift_backend_go/notification-service/internal/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func Audit(logger *pkg.Logger) grpc.UnaryServerInterceptor {
	if logger == nil {
		logger = pkg.NewJSONLogger(pkg.LevelInfo)
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		st := status.Convert(err)
		logger.Info("grpc_audit",
			pkg.LogField("method", info.FullMethod),
			pkg.LogField("code", st.Code().String()),
			pkg.LogField("latency_ms", time.Since(start).Milliseconds()),
		)
		return resp, err
	}
}
