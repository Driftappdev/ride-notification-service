package httpintegration

import "context"

type Client interface {
	Do(ctx context.Context, method, url string, body []byte, headers map[string]string) (int, []byte, error)
}
