package handler

import "context"

type RequestHandlerPort interface {
	Handle(ctx context.Context, payload []byte) error
}
