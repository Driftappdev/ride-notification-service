package servicecore

// CapabilityProfile documents enterprise package targets for this service.
type CapabilityProfile struct {
	ServiceName         string
	InboundHTTP         bool
	InboundGRPC         bool
	InboundEvent        bool
	OutboundGRPCClient  bool
	OutboundEvent       bool
	HasDBOrMigration    bool
	RecommendedPackages []string
}

func EnterpriseCapabilityProfile() CapabilityProfile {
	return CapabilityProfile{
		ServiceName:         "notification-service",
		InboundHTTP:         true,
		InboundGRPC:         false,
		InboundEvent:        true,
		OutboundGRPCClient:  true,
		OutboundEvent:       true,
		HasDBOrMigration:    false,
		RecommendedPackages: []string{"engine-core/messaging", "engine-core/messaging/idempotency", "engine-core/runtime/app", "engine-core/transport/grpc", "middleware/auth", "middleware/metrics", "middleware/ratelimit", "middleware/recovery", "middleware/requestid", "middleware/retry", "middleware/timeout", "middleware/tracing", "middleware/validation"},
	}
}
