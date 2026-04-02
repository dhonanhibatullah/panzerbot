package adaptersoutperipheralrtcmedia

import (
	"context"
	"sync"

	config "github.com/dhonanhibatullah/panzerbot/backend/internal/config"
	domainmodel "github.com/dhonanhibatullah/panzerbot/backend/internal/domain/model"
	portsoutlogging "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/logging"
	portsoutperipheral "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/peripheral"
	"github.com/dhonanhibatullah/panzerbot/backend/internal/utils"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/pion/mediadevices"
	"github.com/pion/mediadevices/pkg/codec/opus"
	"github.com/pion/mediadevices/pkg/codec/vpx"
	"github.com/pion/mediadevices/pkg/prop"
	"github.com/pion/webrtc/v4"

	// Register ALSA and V4L2 drivers with mediadevices.
	_ "github.com/pion/mediadevices/pkg/driver/camera"
	_ "github.com/pion/mediadevices/pkg/driver/microphone"
)

const path = "adapters/out/peripheral/rtc_media"

type peerEntry struct {
	pc            *webrtc.PeerConnection
	cancelSpeaker context.CancelFunc
}

type rtcMedia struct {
	log         portsoutlogging.Log
	api         *webrtc.API
	mediaStream mediadevices.MediaStream
	speakerRate beep.SampleRate
	peers       map[string]*peerEntry
	mu          sync.Mutex
}

func New(
	log portsoutlogging.Log,
	speakerRate beep.SampleRate,
) (
	portsoutperipheral.RTCMedia,
	error,
) {
	opusParams, err := opus.NewParams()
	if err != nil {
		return nil, err
	}
	vpxParams, err := vpx.NewVP8Params()
	if err != nil {
		return nil, err
	}
	vpxParams.BitRate = config.RTCVideoBitRate

	codecSelector := mediadevices.NewCodecSelector(
		mediadevices.WithAudioEncoders(&opusParams),
		mediadevices.WithVideoEncoders(&vpxParams),
	)

	mediaEngine := webrtc.MediaEngine{}
	codecSelector.Populate(&mediaEngine)
	api := webrtc.NewAPI(webrtc.WithMediaEngine(&mediaEngine))

	stream, err := mediadevices.GetUserMedia(mediadevices.MediaStreamConstraints{
		Audio: func(c *mediadevices.MediaTrackConstraints) {
			c.SampleRate = prop.Int(int(config.RTCOpusSampleRate))
			c.ChannelCount = prop.Int(1)
		},
		Video: func(c *mediadevices.MediaTrackConstraints) {
			c.DeviceID = prop.String(config.RTCVideoDevice)
		},
		Codec: codecSelector,
	})
	if err != nil {
		return nil, err
	}

	return &rtcMedia{
		log:         log,
		api:         api,
		mediaStream: stream,
		speakerRate: speakerRate,
		peers:       make(map[string]*peerEntry),
	}, nil
}

func (r *rtcMedia) CreatePeer(ctx context.Context, peerID string, onICECandidate func(domainmodel.ICECandidateData)) (string, error) {
	const tag = path + "/CreatePeer"

	pc, err := r.api.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{config.RTCStunServer}},
		},
	})
	if err != nil {
		r.log.Error(
			ctx, tag,
			"Failed to create peer connection",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
		return "", err
	}

	pc.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}
		init := c.ToJSON()
		sdpMid := ""
		if init.SDPMid != nil {
			sdpMid = *init.SDPMid
		}
		sdpMLineIndex := uint16(0)
		if init.SDPMLineIndex != nil {
			sdpMLineIndex = *init.SDPMLineIndex
		}
		onICECandidate(domainmodel.ICECandidateData{
			Candidate:     init.Candidate,
			SDPMid:        sdpMid,
			SDPMLineIndex: sdpMLineIndex,
		})
	})

	pc.OnTrack(func(remoteTrack *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		if remoteTrack.Kind() != webrtc.RTPCodecTypeAudio {
			return
		}
		go r.handleRemoteAudio(peerID, remoteTrack)
	})

	for _, track := range r.mediaStream.GetTracks() {
		if _, err := pc.AddTrack(track); err != nil {
			pc.Close()
			r.log.Error(
				ctx, tag,
				"Failed to add track",
				domainmodel.LogMeta{
					"error": err.Error(),
				},
			)
			return "", err
		}
	}

	offer, err := pc.CreateOffer(nil)
	if err != nil {
		pc.Close()
		r.log.Error(
			ctx, tag,
			"Failed to create offer",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
		return "", err
	}
	if err := pc.SetLocalDescription(offer); err != nil {
		pc.Close()
		r.log.Error(
			ctx, tag,
			"Failed to set local description",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
		return "", err
	}

	r.mu.Lock()
	r.peers[peerID] = &peerEntry{pc: pc}
	r.mu.Unlock()

	return offer.SDP, nil
}

func (r *rtcMedia) SetAnswer(ctx context.Context, peerID string, answerSDP string) error {
	const tag = path + "/SetAnswer"

	r.mu.Lock()
	entry, ok := r.peers[peerID]
	r.mu.Unlock()
	if !ok {
		return nil
	}

	if err := entry.pc.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  answerSDP,
	}); err != nil {
		r.log.Error(
			ctx, tag,
			"Failed to set remote description",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
		return err
	}
	return nil
}

func (r *rtcMedia) AddICECandidate(ctx context.Context, peerID string, candidate domainmodel.ICECandidateData) error {
	const tag = path + "/AddICECandidate"

	r.mu.Lock()
	entry, ok := r.peers[peerID]
	r.mu.Unlock()
	if !ok {
		return nil
	}

	sdpMid := candidate.SDPMid
	sdpMLineIndex := candidate.SDPMLineIndex
	if err := entry.pc.AddICECandidate(webrtc.ICECandidateInit{
		Candidate:     candidate.Candidate,
		SDPMid:        &sdpMid,
		SDPMLineIndex: &sdpMLineIndex,
	}); err != nil {
		r.log.Error(
			ctx, tag,
			"Failed to add ICE candidate",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
		return err
	}
	return nil
}

func (r *rtcMedia) ClosePeer(ctx context.Context, peerID string) error {
	const tag = path + "/ClosePeer"

	r.mu.Lock()
	entry, ok := r.peers[peerID]
	if ok {
		delete(r.peers, peerID)
	}
	r.mu.Unlock()

	if !ok {
		return nil
	}
	if entry.cancelSpeaker != nil {
		entry.cancelSpeaker()
	}
	if err := entry.pc.Close(); err != nil {
		r.log.Error(
			ctx, tag,
			"Failed to close peer connection",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
		return err
	}
	return nil
}

func (r *rtcMedia) handleRemoteAudio(peerID string, track *webrtc.TrackRemote) {
	speakerCtx, cancel := context.WithCancel(context.Background())

	r.mu.Lock()
	if entry, ok := r.peers[peerID]; ok {
		entry.cancelSpeaker = cancel
	} else {
		cancel()
		r.mu.Unlock()
		return
	}
	r.mu.Unlock()

	streamer := utils.NewBeepRtcStreamer(
		speakerCtx,
		track,
		config.RTCBeepStreamerBuffer,
		config.RTCOpusFrameBuffer,
		config.RTCOpusSamplesPerChannel,
	)

	var s beep.Streamer = streamer
	if r.speakerRate != config.RTCOpusSampleRate {
		s = beep.Resample(4, config.RTCOpusSampleRate, r.speakerRate, streamer)
	}

	speaker.Play(s)
}
