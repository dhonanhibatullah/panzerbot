package utils

import (
	"context"

	config "github.com/dhonanhibatullah/panzerbot/backend/internal/config"
	pionopus "github.com/pion/opus"
	"github.com/pion/webrtc/v4"
)

type BeepRtcStreamer struct {
	samples chan [2]float64
	ctx     context.Context
}

func NewBeepRtcStreamer(ctx context.Context, track *webrtc.TrackRemote) *BeepRtcStreamer {
	b := &BeepRtcStreamer{
		samples: make(chan [2]float64, config.RTCBeepStreamerBuffer),
		ctx:     ctx,
	}
	go b.readLoop(track)
	return b
}

func (b *BeepRtcStreamer) Stream(out [][2]float64) (n int, ok bool) {
	for n = range out {
		select {
		case sample, open := <-b.samples:
			if !open {
				return n, false
			}
			out[n] = sample
		case <-b.ctx.Done():
			return n, false
		}
	}
	return len(out), true
}

func (b *BeepRtcStreamer) Err() error { return nil }

func (b *BeepRtcStreamer) readLoop(track *webrtc.TrackRemote) {
	dec := pionopus.NewDecoder()
	pcm := make([]float32, config.RTCOpusFrameBuffer)

	for {
		select {
		case <-b.ctx.Done():
			return
		default:
		}

		pkt, _, err := track.ReadRTP()
		if err != nil {
			return
		}

		_, isStereo, err := dec.DecodeFloat32(pkt.Payload, pcm)
		if err != nil {
			continue
		}

		if isStereo {
			for i := range config.RTCOpusSamplesPerChannel {
				l := float64(pcm[i*2])
				r := float64(pcm[i*2+1])
				select {
				case b.samples <- [2]float64{l, r}:
				case <-b.ctx.Done():
					return
				default:
				}
			}
		} else {
			for i := range config.RTCOpusSamplesPerChannel {
				v := float64(pcm[i])
				select {
				case b.samples <- [2]float64{v, v}:
				case <-b.ctx.Done():
					return
				default:
				}
			}
		}
	}
}
