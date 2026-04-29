package pkg

import (
	"context"
	goauth "github.com/driftappdev/libpackage/goauth"
)

type Claims = goauth.Claims
type AuthManager = goauth.Manager
type AuthConfig = goauth.Config

func NewAuthManager(cfg AuthConfig) *AuthManager            { return goauth.NewManager(cfg) }
func ClaimsFromContext(ctx context.Context) (*Claims, bool) { return goauth.FromContext(ctx) }
