package middleware

import (
	"context"
	"net/http"

	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/config"
	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/logger"
	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/utils"
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
		httpRequestCounter.Add(context.Background(), 1, metric.WithAttributes(attrs...))
		logger.Log.Debug("Sum 1 to %s metric", "http_request_total")
		next.ServeHTTP(w, r)
	})
}

func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		r = r.WithContext(ctx)

		traceId := utils.GetTraceId(ctx)

		w.Header().Set("X-Trace-ID", traceId)
		// Call the next handler, which can be another middleware function or the final handler
		next.ServeHTTP(w, r)
	})
}

// CORSHeadersMiddleware sets CORS headers for each HTTP request
func CORSHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin (consider specifying the exact origins for better security)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		// If it's a pre-flight request, respond immediately
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Otherwise, pass to the next middleware or final handler
		next.ServeHTTP(w, r)
	})
}
