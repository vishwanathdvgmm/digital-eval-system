package core

import (
	"context"
)

type ctxKey string

const (
	// TraceIDKey is the context key for a per-request trace id.
	TraceIDKey ctxKey = "trace_id"
)

// WithTraceID returns a context with the provided trace id.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// TraceIDFromContext extracts trace id; empty string if absent.
func TraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(TraceIDKey).(string); ok {
		return v
	}
	return ""
}
