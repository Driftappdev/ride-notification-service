package httpport

import "context"

type HTTPHandlerPort interface {
	Handle(ctx context.Context, payload []byte) ([]byte, error)
}
