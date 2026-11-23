package gokit

import (
	"io"
	"log/slog"
	"time"
)

const datetimeMilli = time.DateTime + ".000"

func NewLogger(w io.Writer, level Level) *slog.Logger {
	jsonHandler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				t := a.Value.Time()
				return slog.String(slog.TimeKey, t.Format(datetimeMilli))
			}
			return a
		},
	})

	return slog.New(jsonHandler)
}

type Level = slog.Level

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

func SetLogger(w io.Writer, level Level) {
	l := NewLogger(w, level)
	slog.SetDefault(l)
}

func LogDebug(msg string, args ...any) {
	slog.Default().Debug(msg, args...)
}

func LogInfo(msg string, args ...any) {
	slog.Default().Info(msg, args...)
}

func LogWarn(msg string, args ...any) {
	slog.Default().Warn(msg, args...)
}

func LogError(msg string, args ...any) {
	slog.Default().With("caller", GetCaller(0)).Error(msg, args...)
}
