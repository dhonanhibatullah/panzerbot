package domainappcore

import (
	"context"
	"log/slog"
	"os"

	adaptersoutlogginggeneric "github.com/dhonanhibatullah/panzerbot/backend/internal/adapters/out/logging/generic"
	"github.com/dhonanhibatullah/panzerbot/backend/internal/config"
	domainmodel "github.com/dhonanhibatullah/panzerbot/backend/internal/domain/model"
	"github.com/dhonanhibatullah/panzerbot/backend/sound"
	"github.com/gin-gonic/gin"
	"github.com/peergum/go-rpio/v5"
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

	err = rpio.Open()
	if err != nil {
		adaptersoutlogginggeneric.New(logger).Error(
			ctx, tag,
			"Failed to open RPIO",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
	}

	rpio.StartPwm()
	config.MotorRightAPin.Output()
	config.MotorRightBPin.Output()
	config.MotorRightPwmPin.Mode(rpio.Pwm)
	config.MotorRightPwmPin.Freq(config.PwmFrequency)
	config.MotorLeftAPin.Output()
	config.MotorLeftBPin.Output()
	config.MotorLeftPwmPin.Mode(rpio.Pwm)
	config.MotorLeftPwmPin.Freq(config.PwmFrequency)
	config.ServoPanPin.Output()
	config.ServoTiltPin.Output()

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
	ginEngine.Use(gin.Recovery())

	c.infrastructure = &Infrastructure{
		logger:      logger,
		soundTracks: soundTracks,
		ginEngine:   ginEngine,
	}
	return nil
}

func (c *Core) InfrastructureCleanup() {
	rpio.StopPwm()
	rpio.Close()
}
