package repositoryport

import "context"

type RepositoryPort interface {
	Ping(ctx context.Context) error
}
