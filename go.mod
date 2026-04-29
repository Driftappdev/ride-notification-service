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
	github.com/nats-io/nats.go v1.47.0 // indirect
	github.com/nats-io/nkeys v0.4.11 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260120221211-b8f7ae30c516 // indirect
)

require (
	github.com/driftappdev/libpackage/goauth v0.0.0
	github.com/driftappdev/libpackage/gocircuit v0.0.0
	github.com/driftappdev/libpackage/goerror v0.0.0
	github.com/driftappdev/libpackage/gologger v0.0.0
	github.com/driftappdev/libpackage/gometrics v0.0.0
	github.com/driftappdev/libpackage/goratelimit v0.0.0
	github.com/driftappdev/libpackage/goretry v0.0.0
	github.com/driftappdev/libpackage/gosanitizer v0.0.0
	github.com/driftappdev/libpackage/gotimeout v0.0.0
	github.com/driftappdev/libpackage/gotracing v0.0.0
	github.com/driftappdev/libpackage/logmid/logging-middleware v0.0.0
	github.com/driftappdev/libpackage/resilience/cache v0.0.0
	github.com/driftappdev/libpackage/resilience/pagination v0.0.0
	github.com/driftappdev/libpackage/resilience/validate v0.0.0
	github.com/driftappdev/libpackage/resilience/validator v0.0.0
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
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
)

replace github.com/driftappdev/libpackage/goauth => ../../libpackage/goauth

replace github.com/driftappdev/libpackage/resilience/cache => ../../libpackage/resilience/cache

replace github.com/driftappdev/libpackage/gocircuit => ../../libpackage/gocircuit

replace github.com/driftappdev/libpackage/goerror => ../../libpackage/goerror

replace github.com/driftappdev/libpackage/gologger => ../../libpackage/gologger

replace github.com/driftappdev/libpackage/logmid/logging-middleware => ../../libpackage/logging-middleware

replace github.com/driftappdev/libpackage/gometrics => ../../libpackage/gometrics

replace github.com/driftappdev/libpackage/resilience/pagination => ../../libpackage/resilience/pagination

replace github.com/driftappdev/libpackage/goratelimit => ../../libpackage/goratelimit

replace github.com/driftappdev/libpackage/goretry => ../../libpackage/goretry

replace github.com/driftappdev/libpackage/gosanitizer => ../../libpackage/gosanitizer

replace github.com/driftappdev/libpackage/gotimeout => ../../libpackage/gotimeout

replace github.com/driftappdev/libpackage/gotracing => ../../libpackage/gotracing

replace github.com/driftappdev/libpackage/resilience/validate => ../../libpackage/resilience/validate

replace github.com/driftappdev/libpackage/resilience/validator => ../../libpackage/resilience/validator
