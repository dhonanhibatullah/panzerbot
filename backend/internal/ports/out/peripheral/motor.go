package portsoutperipheral

import "context"

type Motor interface {
	SetSpeedScale(ctx context.Context, scale float64) (err error)
}
