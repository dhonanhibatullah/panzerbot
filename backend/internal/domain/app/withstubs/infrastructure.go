package domainappwithstubs

import (
	"context"
	"log/slog"
	"os"

	adaptersoutlogginggeneric "github.com/dhonanhibatullah/panzerbot/backend/internal/adapters/out/logging/generic"
	"github.com/dhonanhibatullah/panzerbot/backend/internal/config"
	domainmodel "github.com/dhonanhibatullah/panzerbot/backend/internal/domain/model"
	"github.com/dhonanhibatullah/panzerbot/backend/sound"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Infrastructure struct {
	logger      *slog.Logger
	soundTracks []sound.Tracks
	ginEngine   *gin.Engine
}

func (c *Core) NewInfrastructure(ctx context.Context) (err error) {
	const tag = path + "/NewInfrastructure"

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: config.LogLevel,
		}),
	)

	err = config.LoadEnv()
	if err != nil {
		adaptersoutlogginggeneric.New(logger).Error(
			ctx, tag,
			"Failed to load .env",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
		return err
	}

	soundTracks, err := sound.LoadTracks()
	if err != nil {
		adaptersoutlogginggeneric.New(logger).Error(
			ctx, tag,
			"Failed to load soundtracks",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
		return err
	}

	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	_ = ginEngine.SetTrustedProxies(nil)
	ginEngine.Use(
		gin.Recovery(),
		cors.New(cors.Config{
			AllowOrigins:     config.CorsAllowedOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowWebSockets:  true,
			AllowCredentials: true,
		}))

	c.infrastructure = &Infrastructure{
		logger:      logger,
		soundTracks: soundTracks,
		ginEngine:   ginEngine,
	}
	return nil
}

func (c *Core) InfrastructureCleanup() {
}
