package toolkit

import (
	"io"
	"log/slog"
	"os"
)

func init() {
	SetDefaultLog(os.Stdout, slog.LevelDebug)
}

func SetDefaultLog(w io.Writer, level slog.Level) {
	jsonHandler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				t := a.Value.Time()
				return slog.String(slog.TimeKey, t.Format(DatetimeMilli))
			}
			return a
		},
	})

	slog.SetDefault(slog.New(jsonHandler))
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
	// 捕获链路
	stack := GetStack()
	slog.Default().With("stack", stack).Error(msg, args...)
}
