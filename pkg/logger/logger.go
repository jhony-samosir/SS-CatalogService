package logger

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

type contextKey string

const CorrelationIDKey contextKey = "correlation_id"

// ContextHandler enriches slog records with the correlation_id from the context.
type ContextHandler struct {
	slog.Handler
}

// NewContextHandler creates a new ContextHandler wrapping the target handler.
func NewContextHandler(target slog.Handler) *ContextHandler {
	return &ContextHandler{Handler: target}
}

// Handle extracts correlation_id from context and adds it to the log record.
func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx != nil {
		if cid, ok := ctx.Value(CorrelationIDKey).(string); ok && cid != "" {
			r.AddAttrs(slog.String("correlation_id", cid))
		}
		
		// Inject OTel trace_id and span_id from active span
		if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
			sc := span.SpanContext()
			r.AddAttrs(
				slog.String("trace_id", sc.TraceID().String()),
				slog.String("span_id", sc.SpanID().String()),
			)
		}
	}
	return h.Handler.Handle(ctx, r)
}

// SetupLogger initializes the global default slog logger with JSON formatting and ContextHandler.
func SetupLogger() *slog.Logger {
	addSource := false
	if os.Getenv("LOG_ADD_SOURCE") == "true" {
		addSource = true
	}

	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: addSource,
	})

	handler := NewContextHandler(baseHandler)
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}
