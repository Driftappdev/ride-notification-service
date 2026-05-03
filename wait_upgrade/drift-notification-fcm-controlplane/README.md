# drift notification fcm controlplane

Production-grade Go notification core for:
- FCM push notification sender
- notification domain model
- control plane routing/policy/idempotency/rate-limit orchestration

Import path examples:

```go
github.com/driftappdev/libpackage/notification/domain
github.com/driftappdev/libpackage/notification/fcm
github.com/driftappdev/libpackage/notification/controlplane
```

Firebase note:
- FCM Admin SDK send APIs support token/topic/condition and batch operations.
- Multicast/list sends should be chunked to max 500 tokens/messages per call.
