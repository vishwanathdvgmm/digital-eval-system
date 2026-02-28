package core

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// RequestTracingMiddleware injects a trace id into request context and response header.
func RequestTracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trace := r.Header.Get("X-Trace-Id")
		if trace == "" {
			trace = uuid.NewString()
		}
		ctx := WithTraceID(r.Context(), trace)
		w.Header().Set("X-Trace-Id", trace)
		start := time.Now()
		next.ServeHTTP(w, r.WithContext(ctx))
		logrus.WithFields(logrus.Fields{
			"trace_id": trace,
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": time.Since(start).String(),
		}).Info("request completed")
	})
}

// JSONOnlyMiddleware rejects non-JSON content types for endpoints expecting JSON.
func JSONOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow GET/HEAD without JSON requirement
		if r.Method == http.MethodGet || r.Method == http.MethodHead {
			next.ServeHTTP(w, r)
			return
		}
		ct := r.Header.Get("Content-Type")
		if ct == "" || (ct != "application/json" && ct != "application/json; charset=utf-8") {
			w.Header().Set("Content-Type", "application/json")
			payload, _ := MarshalEnvelope(false, nil, NewErrBody(CodeBadRequest, "Content-Type must be application/json", ct), nil)
			w.WriteHeader(http.StatusUnsupportedMediaType)
			_, _ = w.Write(payload)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// WithTimeoutMiddleware sets a default request timeout via context if not present.
func WithTimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			// if ctx already has a deadline, respect it
			if _, ok := ctx.Deadline(); !ok && timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, timeout)
				defer cancel()
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
