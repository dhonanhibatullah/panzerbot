package adaptersoutperipheralmotor

import (
	"context"

	portsoutperipheral "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/peripheral"
)

type motorStub struct {
}

func NewStub() portsoutperipheral.Motor {
	return &motorStub{}
}

func (m *motorStub) SetSpeedScale(ctx context.Context, scale float64) (err error) {
	return nil
}
