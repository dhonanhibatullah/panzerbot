package adaptersoutperipheralmotor

import (
	"context"

	portsoutlogging "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/logging"
	portsoutperipheral "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/peripheral"
	"github.com/peergum/go-rpio/v5"
)

const path = "adapters/out/peripheral/motor"

type motor struct {
	log         portsoutlogging.Log
	aPin        *rpio.Pin
	bPin        *rpio.Pin
	pwmPin      *rpio.Pin
	cycleLength uint32
}

func New(
	log portsoutlogging.Log,
	aPin *rpio.Pin,
	bPin *rpio.Pin,
	pwmPin *rpio.Pin,
	cycleLength uint32,
) portsoutperipheral.Motor {
	return &motor{
		log:         log,
		aPin:        aPin,
		bPin:        bPin,
		pwmPin:      pwmPin,
		cycleLength: cycleLength,
	}
}

func (m *motor) SetSpeedScale(ctx context.Context, scale float64) (err error) {
	const tag = path + "/SetSpeedScale"

	if scale > 0.0 {
		m.aPin.High()
		m.bPin.Low()
	} else if scale < 0.0 {
		m.aPin.Low()
		m.bPin.High()
	} else {
		m.aPin.Low()
		m.bPin.Low()
		m.pwmPin.DutyCycle(0, m.cycleLength)
		return nil
	}

	if scale > 1.0 {
		scale = 1.0
	} else if scale < -1.0 {
		scale = -1.0
	}

	dutyLength := uint32(scale * float64(m.cycleLength))
	m.pwmPin.DutyCycle(dutyLength, m.cycleLength)
	return nil
}
