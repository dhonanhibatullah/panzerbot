package utils

/*
#cgo pkg-config: opus
#include <opus.h>
#include <stdlib.h>
*/
import "C"
import (
	"context"
	"fmt"
	"unsafe"

	"github.com/pion/webrtc/v4"
)

const opusDecodeSampleRate = 48000

type opusDecoder struct {
	dec *C.OpusDecoder
}

func newOpusDecoder() (*opusDecoder, error) {
	var code C.int
	dec := C.opus_decoder_create(C.opus_int32(opusDecodeSampleRate), 1, &code)
	if code != C.OPUS_OK {
		return nil, fmt.Errorf("opus_decoder_create: %d", code)
	}
	return &opusDecoder{dec: dec}, nil
}

func (d *opusDecoder) decodeFloat32(data []byte, pcm []float32) (int, error) {
	n := C.opus_decode_float(
		d.dec,
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.opus_int32(len(data)),
		(*C.float)(unsafe.Pointer(&pcm[0])),
		C.int(len(pcm)),
		0,
	)
	if n < 0 {
		return 0, fmt.Errorf("opus_decode_float error: %d", n)
	}
	return int(n), nil
}

func (d *opusDecoder) destroy() {
	C.opus_decoder_destroy(d.dec)
}

type BeepRtcStreamer struct {
	samples chan [2]float64
	ctx     context.Context
}

func NewBeepRtcStreamer(
	ctx context.Context,
	track *webrtc.TrackRemote,
	streamerBufferSize int,
	opusFrameBuffer int,
	opusSamplePerChannel int,
) *BeepRtcStreamer {
	b := &BeepRtcStreamer{
		samples: make(chan [2]float64, streamerBufferSize),
		ctx:     ctx,
	}
	go b.readLoop(track, opusFrameBuffer)
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
		default:
			out[n] = [2]float64{}
		}
	}
	return len(out), true
}

func (b *BeepRtcStreamer) Err() error { return nil }

func (b *BeepRtcStreamer) readLoop(track *webrtc.TrackRemote, opusFrameBuffer int) {
	dec, err := newOpusDecoder()
	if err != nil {
		return
	}
	defer dec.destroy()

	pcm := make([]float32, opusFrameBuffer)

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

		n, err := dec.decodeFloat32(pkt.Payload, pcm)
		if err != nil || n == 0 {
			continue
		}

		for i := range n {
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
