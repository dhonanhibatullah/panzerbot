package domainappwithstubs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dhonanhibatullah/panzerbot/backend/internal/config"
	domainmodel "github.com/dhonanhibatullah/panzerbot/backend/internal/domain/model"
)

const path = "domain/app/core"

type Core struct {
	infrastructure *Infrastructure
	wiring         *Wiring
	delivery       *Delivery
}

func Run() int {
	const tag = path + "/Run"

	c := Core{}
	ctx := context.Background()

	err := c.NewInfrastructure(ctx)
	if err != nil {
		return 1
	}
	defer c.InfrastructureCleanup()

	err = c.NewWiring(ctx)
	if err != nil {
		return 1
	}

	err = c.NewDelivery(ctx)
	if err != nil {
		return 1
	}

	addr := fmt.Sprintf("0.0.0.0:%d", config.AppPort)
	c.wiring.logPort.Info(
		ctx, tag,
		"Running server...",
		domainmodel.LogMeta{
			"port": config.AppPort,
		},
	)

	err = c.infrastructure.ginEngine.Run(addr)
	if err != nil && err != http.ErrServerClosed {
		c.wiring.logPort.Error(
			ctx, tag,
			"Failed to start HTTP server",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
	}

	return 0
}
