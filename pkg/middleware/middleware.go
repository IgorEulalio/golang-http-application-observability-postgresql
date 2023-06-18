package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/logger"
	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/metrics"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var serviceName = os.Getenv("SERVICE_NAME")
var httpRequestCounter metric.Int64Counter

func init() {
	err := InitHTTPRequestCounter()
	if err != nil {
		// Here you can't return the error so you might want to handle it accordingly
		logger.Log.Fatalf("Failed to initialize HTTP request counter: %v", err)
	}
}

func InitHTTPRequestCounter() error {
	ctx := context.Background()
	meterProvider, err := metrics.InitMetricsProvider(ctx)
	if err != nil {
		return err
	}

	meter := meterProvider.Meter("repositories-service")

	httpRequestCounter, err = meter.Int64Counter(
		"http_request_total",
		metric.WithDescription("Counts total HTTP requests"),
	)
	if err != nil {
		return err
	}

	return nil
}

func HTTPRequestCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if httpRequestCounter == nil {
			logger.Log.Error("httpRequestCounter is not initialized")
			return
		}

		attrs := []attribute.KeyValue{
			attribute.String("method", r.Method),
			attribute.String("path", r.URL.Path),
			attribute.String("service", serviceName),
		}
		logger.Log.Debug("metric added...")
		httpRequestCounter.Add(context.Background(), 1, metric.WithAttributes(attrs...))
		next.ServeHTTP(w, r)
	})
}
