module github.com/IgorEulalio/golang-http-application-observability-postgresql

go 1.20

require (
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/google/uuid v1.3.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.10.9
	go.opentelemetry.io/otel v1.16.0
	go.opentelemetry.io/otel/metric v1.16.0
	go.opentelemetry.io/otel/sdk v1.16.0
	go.opentelemetry.io/otel/sdk/metric v0.39.0
	go.opentelemetry.io/otel/trace v1.16.0
)

require (
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.0 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.16.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.39.0 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
)

require (
	github.com/gorilla/mux v1.8.0
	github.com/sirupsen/logrus v1.9.2
	github.com/streadway/amqp v1.1.0
	github.com/uptrace/opentelemetry-go-extra/otelsql v0.2.2
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.42.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.42.0
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.39.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.16.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.16.0
	golang.org/x/sys v0.8.0 // indirect
	google.golang.org/grpc v1.56.3
)
