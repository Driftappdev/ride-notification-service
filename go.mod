module dift_backend_go/notification-service

go 1.25.1

require (
	github.com/eclipse/paho.mqtt.golang v1.3.5
	github.com/segmentio/kafka-go v0.4.49
	github.com/spf13/viper v1.21.0
	google.golang.org/grpc v1.80.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/nats-io/nats.go v1.47.0 // indirect
	github.com/nats-io/nkeys v0.4.11 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	golang.org/x/crypto v0.51.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260401024825-9d38bb4040a9 // indirect
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect; indirect, required by viper
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/sagikazarmark/locafero v0.11.0 // indirect
	github.com/sourcegraph/conc v0.3.1-0.20240121214520-5f936abd7ae8 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
)

require github.com/PlatformCore/middleware v0.0.0

require github.com/PlatformCore/engine-core/runtime v0.0.0

require github.com/PlatformCore/engine-core/messaging v0.0.0

require github.com/PlatformCore/engine-core/observability v0.0.0

require github.com/PlatformCore/engine-core/resilience v0.0.0

replace github.com/PlatformCore/middleware => ../middleware

replace github.com/PlatformCore/engine-core/runtime => ../engine-core/runtime

replace github.com/PlatformCore/engine-core/messaging => ../engine-core/messaging

replace github.com/PlatformCore/engine-core/observability => ../engine-core/observability

replace github.com/PlatformCore/engine-core/resilience => ../engine-core/resilience

require github.com/PlatformCore/engine-core/transport v0.0.0

require github.com/PlatformCore/engine-core/security v0.0.0

require github.com/PlatformCore/engine-core/validation v0.0.0

require github.com/PlatformCore/engine-core/tenant v0.0.0

require github.com/PlatformCore/engine-core/plugins v0.0.0

replace github.com/PlatformCore/engine-core/transport => ../engine-core/transport

replace github.com/PlatformCore/engine-core/security => ../engine-core/security

replace github.com/PlatformCore/engine-core/validation => ../engine-core/validation

replace github.com/PlatformCore/engine-core/tenant => ../engine-core/tenant

replace github.com/PlatformCore/engine-core/plugins => ../engine-core/plugins
