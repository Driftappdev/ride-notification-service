package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Permission() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			if vals := md.Get("x-internal-auth"); len(vals) > 0 && vals[0] == "deny" {
				return nil, status.Error(codes.PermissionDenied, "internal permission denied")
			}
		}
		return handler(ctx, req)
	}
}
