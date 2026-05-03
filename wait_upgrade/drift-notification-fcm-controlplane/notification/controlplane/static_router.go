package controlplane

import (
	"context"

	"github.com/driftappdev/libpackage/notification/domain"
	"github.com/driftappdev/libpackage/notification/port"
)

type StaticRouter struct {
	routes map[domain.Channel]port.Sender
}

func NewStaticRouter(routes map[domain.Channel]port.Sender) *StaticRouter {
	return &StaticRouter{routes: routes}
}

func (r *StaticRouter) Route(ctx context.Context, n domain.Notification) (port.Sender, error) {
	s, ok := r.routes[n.Channel]
	if !ok || s == nil {
		return nil, domain.ErrNoRoute
	}
	return s, nil
}
