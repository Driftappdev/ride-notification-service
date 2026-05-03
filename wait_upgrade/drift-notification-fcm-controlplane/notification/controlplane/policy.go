package controlplane

import (
	"context"
	"time"

	"github.com/driftappdev/libpackage/notification/domain"
	"github.com/driftappdev/libpackage/notification/port"
)

type QuietHoursPolicy struct {
	StartHour int
	EndHour   int
	Location  *time.Location
}

func (p QuietHoursPolicy) Evaluate(ctx context.Context, n domain.Notification) (port.PolicyDecision, error) {
	if n.Priority >= domain.PriorityCritical {
		return port.PolicyDecision{Allow: true}, nil
	}
	loc := p.Location
	if loc == nil { loc = time.Local }
	h := time.Now().In(loc).Hour()

	inQuiet := false
	if p.StartHour < p.EndHour {
		inQuiet = h >= p.StartHour && h < p.EndHour
	} else {
		inQuiet = h >= p.StartHour || h < p.EndHour
	}
	if inQuiet {
		return port.PolicyDecision{Allow: false, Reason: "quiet_hours"}, nil
	}
	return port.PolicyDecision{Allow: true}, nil
}

type ChainPolicy struct {
	Policies []port.PolicyEngine
}

func (c ChainPolicy) Evaluate(ctx context.Context, n domain.Notification) (port.PolicyDecision, error) {
	for _, p := range c.Policies {
		if p == nil { continue }
		d, err := p.Evaluate(ctx, n)
		if err != nil { return d, err }
		if d.Mutate != nil { d.Mutate(&n) }
		if !d.Allow { return d, nil }
	}
	return port.PolicyDecision{Allow: true}, nil
}
