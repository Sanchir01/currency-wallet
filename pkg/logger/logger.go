package logger

import (
	"log/slog"
	"os"
)

var (
	development = "development"
	production  = "production"
)

func SetupLogger(env string) *slog.Logger {
	var lg *slog.Logger
	switch env {
	case production:
		lg = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case development:
		lg = setupPrettySlog()
	}
	return lg
}

func setupPrettySlog() *slog.Logger {
	opts := PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
