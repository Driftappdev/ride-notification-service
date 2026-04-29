package middleware

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var idempotencyStore sync.Map

func Idempotency() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}
		keys := md.Get("x-idempotency-key")
		if len(keys) == 0 {
			return handler(ctx, req)
		}

		k := info.FullMethod + ":" + keys[0]
		if _, exists := idempotencyStore.LoadOrStore(k, time.Now().Unix()); exists {
			return nil, status.Error(codes.AlreadyExists, "duplicate idempotent request")
		}

		resp, err := handler(ctx, req)
		if err != nil {
			idempotencyStore.Delete(k)
		}
		return resp, err
	}
}
