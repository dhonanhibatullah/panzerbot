package adaptersoutperipheralservo

import (
	"context"
	"math"
	"sync"
	"time"

	portsoutlogging "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/logging"
	portsoutperipheral "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/peripheral"
	"github.com/peergum/go-rpio/v5"
)

const path = "adapters/out/peripheral/servo"

const (
	minPulse   = 500 * time.Microsecond
	maxPulse   = 2500 * time.Microsecond
	pulseRange = maxPulse - minPulse
)

type servo struct {
	log          portsoutlogging.Log
	pin          *rpio.Pin
	highDuration time.Duration
	ticker       *time.Ticker
	mu           sync.Mutex
}

func New(
	log portsoutlogging.Log,
	pin *rpio.Pin,
) portsoutperipheral.Servo {
	s := &servo{
		log:          log,
		pin:          pin,
		highDuration: (minPulse + (pulseRange)/2) * time.Microsecond,
		ticker:       time.NewTicker(20 * time.Millisecond),
		mu:           sync.Mutex{},
	}
	go s.servoWorker()
	return s
}

func (s *servo) SetAngle(ctx context.Context, angle float64) (err error) {
	const tag = path + "/SetAngle"

	if angle < 0 {
		angle = 0
	} else if angle > math.Pi {
		angle = math.Pi
	}

	s.mu.Lock()
	s.highDuration = minPulse + time.Duration((angle/math.Pi)*float64(pulseRange))
	s.mu.Unlock()

	return nil
}

func (s *servo) servoWorker() {
	const tag = path + "/servoWorker"

	for range s.ticker.C {
		s.mu.Lock()
		duration := s.highDuration
		s.mu.Unlock()

		s.pin.High()
		time.Sleep(duration)
		s.pin.Low()
	}
}
