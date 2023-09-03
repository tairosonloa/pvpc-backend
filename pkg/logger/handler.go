package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strings"
)

type CustomTextHandler struct {
	slog.TextHandler
	addSource bool
}

type CustomJSONHandler struct {
	slog.JSONHandler
	addSource bool
}

// NewCustomTextHandler initializes a custom text handler that is based on
// slog.TextHandler but adds to the log record, when calling a log function
// with context, the values of the keys defined in constants.go that are
// present in the given context.
func NewCustomTextHandler(w io.Writer, opts *slog.HandlerOptions) *CustomTextHandler {
	// Disable the slog.TextHandler addSource behavior if set,
	// because it overlaps with the custom addSource behavior
	addSource := opts.AddSource || opts.Level == slog.LevelDebug
	opts.AddSource = false
	return &CustomTextHandler{
		TextHandler: *slog.NewTextHandler(w, opts),
		addSource:   addSource,
	}
}

// NewCustomJSONHandler initializes a custom JSON handler that is based on
// slog.JSONHandler but adds to the log record, when calling a log function
// with context, the values of the keys defined in constants.go that are
// present in the given context.
func NewCustomJSONHandler(w io.Writer, opts *slog.HandlerOptions) *CustomJSONHandler {
	// Disable the slog.JSONHandler addSource behavior if set,
	// because it overlaps with the custom addSource behavior
	addSource := opts.AddSource || opts.Level == slog.LevelDebug
	opts.AddSource = false
	return &CustomJSONHandler{
		JSONHandler: *slog.NewJSONHandler(w, opts),
		addSource:   addSource,
	}
}

// Handle implements slog.Handler interface.
func (h *CustomTextHandler) Handle(ctx context.Context, r slog.Record) error {
	addContextInfoToRecord(ctx, &r)
	if h.addSource {
		addSourceToRecord(&r)
	}
	return h.TextHandler.Handle(ctx, r)
}

// Handle implements slog.Handler interface.
func (h *CustomJSONHandler) Handle(ctx context.Context, r slog.Record) error {
	addContextInfoToRecord(ctx, &r)
	if h.addSource {
		addSourceToRecord(&r)
	}
	return h.JSONHandler.Handle(ctx, r)
}

func addContextInfoToRecord(ctx context.Context, r *slog.Record) {
	reqID := ctx.Value(ContextKeyRequestID)
	if value, ok := reqID.(string); ok {
		r.AddAttrs(slog.String("reqID", value))
	}
	contextErr := ctx.Err()
	if contextErr != nil {
		r.AddAttrs(slog.String("contextErr", contextErr.Error()))
	}
}

func addSourceToRecord(r *slog.Record) {
	_, file, line, ok := runtime.Caller(5)
	if ok {
		fileSplit := strings.Split(file, "/")
		r.AddAttrs(slog.String(slog.SourceKey, fmt.Sprintf("%s/%s:%d", fileSplit[len(fileSplit)-2], fileSplit[len(fileSplit)-1], line)))
	}
}

func buildHandlerOptions(opts *slog.HandlerOptions) *slog.HandlerOptions {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	if opts.Level == nil {
		opts.Level = slog.LevelInfo
	}

	if opts.ReplaceAttr == nil {
		opts.ReplaceAttr = replaceAttr
	}

	return opts
}

// replaceAttr updates the log output format by replacing the value of the
// slog.LevelKey attribute with the string representation of the slog.Level.
func replaceAttr(_groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)
		levelLabel, exists := LevelNames[level]
		if !exists {
			levelLabel = level.String()
		}

		a.Value = slog.StringValue(levelLabel)
	}
	return a
}
