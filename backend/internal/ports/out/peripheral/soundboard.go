package portsoutperipheral

import "context"

type Soundboard interface {
	PlayTrack(ctx context.Context, idx int) (err error)
	StopTracks(ctx context.Context) (err error)
	GetTrackNames(ctx context.Context) (trackNames []string, err error)
	GetTrackFileNames(ctx context.Context) (trackFilenames []string, err error)
	SetTrackVolume(ctx context.Context, volume int) (err error)
}
