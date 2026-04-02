package adaptersoutperipheralsoundboard

import (
	"context"
	"math"
	"sync"
	"time"

	portsoutlogging "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/logging"
	portsoutperipheral "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/peripheral"
	"github.com/dhonanhibatullah/panzerbot/backend/sound"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
)

const path = "adapters/out/peripheral/soundboard"

type soundboard struct {
	log    portsoutlogging.Log
	mu     sync.Mutex
	tracks []sound.Tracks
	volume int
}

func New(
	log portsoutlogging.Log,
	tracks []sound.Tracks,
	speakerRate beep.SampleRate,
) portsoutperipheral.Soundboard {
	s := &soundboard{
		log:    log,
		tracks: tracks,
		volume: 100,
	}
	if len(tracks) > 0 {
		speaker.Init(
			speakerRate,
			speakerRate.N(100*time.Millisecond),
		)
	}
	return s
}

func (s *soundboard) PlayTrack(ctx context.Context, idx int) (err error) {
	const tag = path + "/PlayTrack"

	s.mu.Lock()
	defer s.mu.Unlock()
	if idx < 0 || idx >= len(s.tracks) {
		return nil
	}
	return s.play(s.tracks[idx].Stream)
}

func (s *soundboard) StopTracks(ctx context.Context) (err error) {
	const tag = path + "/StopTracks"

	speaker.Lock()
	speaker.Clear()
	speaker.Unlock()
	return nil
}

func (s *soundboard) GetTrackNames(ctx context.Context) (trackNames []string, err error) {
	const tag = path + "/GetTrackNames"

	for _, t := range s.tracks {
		trackNames = append(trackNames, t.Name)
	}
	return
}

func (s *soundboard) GetTrackFileNames(ctx context.Context) (trackFilenames []string, err error) {
	const tag = path + "/GetTrackFileNames"

	for _, t := range s.tracks {
		trackFilenames = append(trackFilenames, t.FileName)
	}
	return
}

func (s *soundboard) SetTrackVolume(ctx context.Context, volume int) (err error) {
	const tag = path + "/SetTrackVolume"

	s.mu.Lock()
	defer s.mu.Unlock()
	if volume < 0 {
		volume = 0
	} else if volume > 100 {
		volume = 100
	}
	s.volume = volume
	return nil
}

func (s *soundboard) play(stream beep.StreamSeekCloser) error {
	const tag = path + "/play"

	if err := stream.Seek(0); err != nil {
		return err
	}

	var dB float64
	if s.volume <= 0 {
		dB = -math.MaxFloat64
	} else {
		dB = float64(s.volume-100) / 10.0
	}

	vol := &effects.Volume{
		Streamer: stream,
		Base:     2,
		Volume:   dB,
	}

	speaker.Lock()
	speaker.Play(vol)
	speaker.Unlock()
	return nil
}
