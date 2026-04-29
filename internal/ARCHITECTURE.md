# Internal Architecture (notification-service)

This service follows layered design:
- internal/adapter: transport and gateway adapters.
- internal/interface: ports/contracts between layers.
- internal/integration: implementations for external systems.
- internal/service_logic: orchestration and use cases.
- internal/dto and internal/model: boundary payloads and domain entities.
- internal/pkg: all libpackage imports are wrapped here.
- internal/servicecore: middleware, health and service-level concerns.
