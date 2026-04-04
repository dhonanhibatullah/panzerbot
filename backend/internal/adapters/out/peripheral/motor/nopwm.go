package adaptersoutperipheralmotor

import (
	"context"

	portsoutlogging "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/logging"
	portsoutperipheral "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/peripheral"
	"github.com/peergum/go-rpio/v5"
)

type motorNoPwm struct {
	log  portsoutlogging.Log
	aPin *rpio.Pin
	bPin *rpio.Pin
}

func NewNoPwm(
	log portsoutlogging.Log,
	aPin *rpio.Pin,
	bPin *rpio.Pin,
) portsoutperipheral.Motor {
	return &motorNoPwm{
		log:  log,
		aPin: aPin,
		bPin: bPin,
	}
}

func (m *motorNoPwm) SetSpeedScale(ctx context.Context, scale float64) (err error) {
	if scale > 0.0 {
		m.aPin.High()
		m.bPin.Low()
	} else if scale < 0.0 {
		m.aPin.Low()
		m.bPin.High()
	} else {
		m.aPin.Low()
		m.bPin.Low()
	}
	return nil
}
