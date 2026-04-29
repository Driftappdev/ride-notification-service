package servicecore

import "dift_backend_go/notification-service/internal/pkg"

type ServiceError = pkg.BaseError

func NewServiceError(code string, status int, message string) *ServiceError {
	return pkg.NewBaseError(code, status, message)
}
