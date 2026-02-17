package logger

import (
	"context"
	"log/slog"
	"time"

	"github.com/taufiksoleh/backend-mp/pkg/reqctx"
)

func FromContext(ctx context.Context) *slog.Logger {
	id := reqctx.RequestID(ctx)
	if id == "" {
		return slog.Default()
	}
	return slog.Default().With("request_id", id)
}

func Track(ctx context.Context, tag string, args ...any) func(*error) {
	log := FromContext(ctx).With("tag", tag)
	start := time.Now()
	return func(errp *error) {
		attrs := append([]any{"duration", time.Since(start)}, args...)
		if *errp != nil {
			log.Error("fail", append([]any{"err", *errp}, attrs...)...)
		} else {
			log.Info("success", attrs...)
		}
	}
}
