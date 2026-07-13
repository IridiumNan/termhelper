package output

import (
	"log/slog"
	"os"
)

var Console *slog.Logger

func InitConsoleLogger(level slog.Level) {
	consoleLevel := level

	consoleOpts := &slog.HandlerOptions{
		Level: consoleLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format("15:04:05"))
			}

			return a
		},
	}

	Console = slog.New(slog.NewTextHandler(os.Stderr, consoleOpts))
}
