package logger

import (
	"log/slog"
	"strings"
)

// Custom slog levels
const (
	LevelFatal = slog.Level(12)
)

// LevelNames maps slog.Levels to their string representation.
var LevelNames = map[slog.Leveler]string{
	LevelFatal: "FATAL",
}

// ParseLevel parses a string level into a slog.Level.
func ParseLevel(level string) slog.Level {
	switch strings.ToUpper(level) {
	case slog.LevelDebug.String():
		return slog.LevelDebug
	case slog.LevelInfo.String():
		return slog.LevelInfo
	case slog.LevelWarn.String():
		return slog.LevelWarn
	case slog.LevelError.String():
		return slog.LevelError
	case LevelNames[LevelFatal]:
		return LevelFatal
	default:
		return slog.LevelInfo
	}
}
