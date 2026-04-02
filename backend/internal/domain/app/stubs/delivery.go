package domainappstubs

import (
	"context"

	adaptersinhttphandler "github.com/dhonanhibatullah/panzerbot/backend/internal/adapters/in/http/handler"
	adaptersinhttproutes "github.com/dhonanhibatullah/panzerbot/backend/internal/adapters/in/http/routes"
)

type Delivery struct {
	peripheralHttpHdl    *adaptersinhttphandler.Peripheral
	rtcSignallingHttpHdl *adaptersinhttphandler.RTCSignalling
	httpRoutes           *adaptersinhttproutes.Routes
}

func (c *Core) NewDelivery(ctx context.Context) (err error) {
	const tag = path + "/NewDelivery"

	peripheralHttpHdl := adaptersinhttphandler.NewPeripheral(
		c.wiring.logPort,
		c.wiring.peripheralSvcPort,
	)
	rtcSignallingHttpHdl := adaptersinhttphandler.NewRTCSignalling(
		c.wiring.logPort,
		c.wiring.rtcSignallingSvcPort,
	)
	httpRoutes := adaptersinhttproutes.New(
		c.infrastructure.ginEngine.Group(""),
		peripheralHttpHdl,
		rtcSignallingHttpHdl,
	)

	httpRoutes.Route()

	c.delivery = &Delivery{
		peripheralHttpHdl:    peripheralHttpHdl,
		rtcSignallingHttpHdl: rtcSignallingHttpHdl,
		httpRoutes:           httpRoutes,
	}
	return nil
}
