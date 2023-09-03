package logger

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"
)

var handler slog.Handler

// SetDefaultLoggerText creates a new slog.Logger with a custom text
// handler and sets it as the default logger.
//
// The custom text handler is created by calling NewCustomTextHandler.
// This handler will add to the log record, when calling a log function
// with context, the values of the keys defined in constants.go that are
// present in the given context.
func SetDefaultLoggerText(opts *slog.HandlerOptions) {
	handler = NewCustomTextHandler(os.Stdout, buildHandlerOptions(opts))
	slog.SetDefault(slog.New(handler))
}

// SetDefaultLoggerJSON creates a new slog.Logger with a custom text
// handler and sets it as the default logger.
//
// The custom text handler is created by calling NewCustomJSONHandler.
// This handler will add to the log record, when calling a log function
// with context, the values of the keys defined in constants.go that are
// present in the given context.
func SetDefaultLoggerJSON(opts *slog.HandlerOptions) {
	handler = NewCustomJSONHandler(os.Stdout, buildHandlerOptions(opts))
	slog.SetDefault(slog.New(handler))
}

// SetTestLogger creates a new slog.Logger with a CustomTextHandler
// handler that logs at LevelDebug and sets it as the default logger.
func SetTestLogger(output io.Writer) {
	if output == nil {
		output = os.Stderr
	}

	handler = slog.NewTextHandler(output, buildHandlerOptions(&slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(slog.New(handler))
}

// ServerInfoLoggerFromDefault returns a new *log.Logger that is intended to be
// injected as ErrorLog attribute of an http.Server.
//
// It, in essence, calls slog.NewLogLogger with the handler created by
// SetDefaultLoggerText or SetDefaultLoggerJson and the LevelError level.
func ServerErrorLoggerFromDefault() *log.Logger {
	if handler == nil {
		panic("default logger not initialized")
	}
	return slog.NewLogLogger(handler, slog.LevelError)
}

// Debug calls slog.Debug on the default logger.
func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

// DebugContext calls slog.DebugContext on the default logger.
func DebugContext(ctx context.Context, msg string, args ...any) {
	slog.DebugContext(ctx, msg, args...)
}

// Info calls slog.Info on the default logger.
func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

// InfoContext calls slog.InfoContext on the default logger.
func InfoContext(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, args...)
}

// Warn calls slog.Warn on the default logger.
func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

// WarnContext calls slog.WarnContext on the default logger.
func WarnContext(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, args...)
}

// Error calls slog.Error on the default logger.
func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

// ErrorContext calls slog.ErrorContext on the default logger.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}

// Fatal logs at LevelFatal.
func Fatal(msg string, args ...any) {
	slog.Log(context.Background(), LevelFatal, msg, args...)
	os.Exit(1)
}

// FatalContext logs at LevelFatal with the given context.
func FatalContext(ctx context.Context, msg string, args ...any) {
	slog.Log(ctx, LevelFatal, msg, args...)
	os.Exit(1)
}

// Log calls slog.Log on the default logger.
func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	slog.Log(ctx, level, msg, args...)
}

// LogAttrs calls slog.LogAttrs on the default logger.
func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	slog.LogAttrs(ctx, level, msg, attrs...)
}
