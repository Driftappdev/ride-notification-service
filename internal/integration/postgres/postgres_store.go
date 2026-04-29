package postgresintegration

import "context"

type Store interface {
	Ping(ctx context.Context) error
}
