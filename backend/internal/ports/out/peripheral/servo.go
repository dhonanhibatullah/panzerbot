package portsoutperipheral

import "context"

type Servo interface {
	SetAngle(ctx context.Context, angle float64) (err error)
}
