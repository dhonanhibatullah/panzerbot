package adaptersoutlogginggeneric

import (
	"context"
	"log/slog"

	domainmodel "github.com/dhonanhibatullah/panzerbot/backend/internal/domain/model"
	portsoutlog "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/logging"
)

type log struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) portsoutlog.Log {
	return &log{
		logger: logger,
	}
}

func (l *log) Error(ctx context.Context, tag string, message string, meta domainmodel.LogMeta) {
	l.logger.ErrorContext(ctx, message,
		slog.String("tag", tag),
		slog.Any("meta", meta),
	)
}

func (l *log) Warn(ctx context.Context, tag string, message string, meta domainmodel.LogMeta) {
	l.logger.WarnContext(ctx, message,
		slog.String("tag", tag),
		slog.Any("meta", meta),
	)
}

func (l *log) Info(ctx context.Context, tag string, message string, meta domainmodel.LogMeta) {
	l.logger.InfoContext(ctx, message,
		slog.String("tag", tag),
		slog.Any("meta", meta),
	)
}

func (l *log) Debug(ctx context.Context, tag string, message string, meta domainmodel.LogMeta) {
	l.logger.DebugContext(ctx, message,
		slog.String("tag", tag),
		slog.Any("meta", meta),
	)
}
