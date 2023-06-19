package middleware

import (
	"context"
	"net/http"

	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/config"
	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var serviceName = config.Config.ServiceName
var httpRequestCounter metric.Int64Counter

func InitMetrics(meter metric.Meter) error {
	var err error
	httpRequestCounter, err = meter.Int64Counter(
		"http_request_total",
		metric.WithDescription("Counts total HTTP requests"),
	)
	if err != nil {
		logger.Log.Error("Failed initializing metrics.")
		return err
	}

	return nil
}

func HTTPRequestCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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
