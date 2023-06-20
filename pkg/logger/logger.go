package logger

import (
	"bytes"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/config"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

var Log *logrus.Entry

// ResponseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type ResponseWriter struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

// NewResponseWriter creates a new responseWriter.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		body:           &bytes.Buffer{},
	}
}

// WriteHeader saves the status code and writes it to the underlying
// http.ResponseWriter.
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write writes the data to the body and underlying http.ResponseWriter.
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// LogRequestResponse logs the details of the HTTP request and response.
func LogRequestResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rw := NewResponseWriter(w)

		next.ServeHTTP(rw, r)

		// Extract the trace ID from the context.
		span := trace.SpanFromContext(r.Context())

		traceID := span.SpanContext().TraceID().String()

		duration := time.Since(startTime)
		Log.WithFields(logrus.Fields{
			"method":          r.Method,
			"path":            r.URL.Path,
			"duration":        duration,
			"requestHost":     r.Host,
			"requestBody":     r.Body,
			"responseCode":    rw.status,
			"responseBody":    rw.body.String(),
			"requestHeaders":  r.Header,
			"responseHeaders": rw.Header(),
			"traceId":         traceID,
		}).Info("Handled request")
	})
}

func init() {

	if config.Config == nil {
		config.LoadConfig()
	}

	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})

	serviceName := config.Config.ServiceName

	logLevel := config.Config.LogLevel

	level, err := logrus.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		logger.Warningf("Invalid or empty log level provided, defaulting to 'info'. Error: %s", err.Error())
		level = logrus.InfoLevel
	}

	logger.SetLevel(level)

	Log = logger.WithField("service", serviceName)
}
