package portsinhttp

import "context"

type Peripheral interface {
	SetMotorSpeedScale(ctx context.Context, right float64, left float64) (err error)
	SetServoAngle(ctx context.Context, pan float64, tilt float64) (err error)
	PlaySoundboardTrack(ctx context.Context, idx int) (err error)
	StopSoundboardTracks(ctx context.Context) (err error)
	GetSoundboardTrackNames(ctx context.Context) (trackNames []string, err error)
	GetSoundboardTrackFileNames(ctx context.Context) (trackFileNames []string, err error)
}
