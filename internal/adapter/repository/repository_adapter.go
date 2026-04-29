package repositoryadapter

import "context"

type RepositoryAdapter interface {
	Ping(ctx context.Context) error
}
