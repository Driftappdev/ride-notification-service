package service

import "context"

type UseCasePort interface {
	Execute(ctx context.Context, input map[string]interface{}) (map[string]interface{}, error)
}
