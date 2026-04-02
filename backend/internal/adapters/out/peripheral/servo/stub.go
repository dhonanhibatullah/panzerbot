package adaptersoutperipheralservo

import (
	"context"

	portsoutperipheral "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/peripheral"
)

type servoStub struct {
}

func NewStub() portsoutperipheral.Servo {
	return &servoStub{}
}

func (s *servoStub) SetAngle(ctx context.Context, angle float64) (err error) {
	return nil
}
