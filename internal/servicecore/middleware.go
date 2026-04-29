package servicecore

import (
	"dift_backend_go/notification-service/internal/pkg"
	coremid "dift_backend_go/notification-service/internal/servicecore/middleware"

	"google.golang.org/grpc"
)

func DefaultUnaryInterceptors(logger *pkg.Logger) []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		coremid.Context(),
		coremid.Audit(logger),
		coremid.Permission(),
		coremid.Idempotency(),
		coremid.Validation(),
	}
}
