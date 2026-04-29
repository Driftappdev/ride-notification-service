package serviceport

import "context"

type ServicePort interface {
	Name() string
	Health(ctx context.Context) error
}
