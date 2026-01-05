package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

func New(out io.Writer, component string) *slog.Logger {
	if out == nil {
		out = os.Stdout
	}

	handler := slog.NewJSONHandler(out, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)
	if component != "" {
		logger = logger.With("component", component)
	}
	return logger
}

func Stdout(component string) *slog.Logger {
	return New(os.Stdout, component)
}

type ctxKey struct{}

func WithContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
