package domain

import "errors"

var (
	ErrDuplicate       = errors.New("duplicate notification")
	ErrSuppressed      = errors.New("notification suppressed by policy")
	ErrRateLimited     = errors.New("notification rate limited")
	ErrNoRoute         = errors.New("no route available")
	ErrInvalidMessage  = errors.New("invalid notification")
	ErrProviderFailure = errors.New("provider failure")
)
