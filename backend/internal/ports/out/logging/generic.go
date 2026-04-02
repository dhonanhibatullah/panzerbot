package portsoutlogging

import (
	"context"

	domainmodel "github.com/dhonanhibatullah/panzerbot/backend/internal/domain/model"
)

type Log interface {
	Error(ctx context.Context, tag string, message string, meta domainmodel.LogMeta)
	Warn(ctx context.Context, tag string, message string, meta domainmodel.LogMeta)
	Info(ctx context.Context, tag string, message string, meta domainmodel.LogMeta)
	Debug(ctx context.Context, tag string, message string, meta domainmodel.LogMeta)
}
