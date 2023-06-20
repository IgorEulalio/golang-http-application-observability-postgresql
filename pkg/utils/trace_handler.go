package utils

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func GetTraceId(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	traceId := span.SpanContext().TraceID().String()

	return traceId
}
