package adaptersoutperipheralservo

import (
	"context"
	"math"

	"github.com/dhonanhibatullah/panzerbot/backend/internal/config"
	portsoutlogging "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/logging"
	portsoutperipheral "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/peripheral"
	"github.com/peergum/go-rpio/v5"
)

type servoPwm struct {
	log portsoutlogging.Log
	pin *rpio.Pin
}

func NewPwm(
	log portsoutlogging.Log,
	pin *rpio.Pin,
) portsoutperipheral.Servo {
	return &servoPwm{
		log: log,
		pin: pin,
	}
}

func (s *servoPwm) SetAngle(ctx context.Context, angle float64) (err error) {
	if angle < 0 {
		angle = 0
	} else if angle > math.Pi {
		angle = math.Pi
	}

	pulseUs := uint32(500 + (angle/math.Pi)*2000)
	s.pin.DutyCycle(pulseUs, config.PwmCycleLength)

	return nil
}
