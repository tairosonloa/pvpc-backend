package logger

import (
	"context"
	"io"
	"log/slog"
)

type CustomTextHandler struct {
	slog.TextHandler
}

type CustomJSONHandler struct {
	slog.JSONHandler
}

// NewCustomTextHandler initializes a custom text handler that is based on
// slog.TextHandler but adds to the log record, when calling a log function
// with context, the values of the keys defined in constants.go that are
// present in the given context.
func NewCustomTextHandler(w io.Writer, opts *slog.HandlerOptions) *CustomTextHandler {
	return &CustomTextHandler{
		TextHandler: *slog.NewTextHandler(w, opts),
	}
}

// NewCustomJSONHandler initializes a custom JSON handler that is based on
// slog.JSONHandler but adds to the log record, when calling a log function
// with context, the values of the keys defined in constants.go that are
// present in the given context.
func NewCustomJSONHandler(w io.Writer, opts *slog.HandlerOptions) *CustomJSONHandler {
	return &CustomJSONHandler{
		JSONHandler: *slog.NewJSONHandler(w, opts),
	}
}

// Handle implements slog.Handler interface.
func (h *CustomTextHandler) Handle(ctx context.Context, r slog.Record) error {
	addContextInfoToRecord(ctx, &r)
	return h.TextHandler.Handle(ctx, r)
}

// Handle implements slog.Handler interface.
func (h *CustomJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	addContextInfoToRecord(ctx, &r)
	return h.JSONHandler.Handle(ctx, r)
}

func addContextInfoToRecord(ctx context.Context, r *slog.Record) {
	reqID := ctx.Value(ContextKeyRequestID)
	if value, ok := reqID.(string); ok {
		r.AddAttrs(slog.String("req_id", value))
	}
}
