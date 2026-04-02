package domainservicehttprtcsignalling

import (
	"context"

	domainmodel "github.com/dhonanhibatullah/panzerbot/backend/internal/domain/model"
	portsinhttp "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/in/http"
	portsoutlogging "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/logging"
	portsoutperipheral "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/peripheral"
)

const path = "domain/service/http/rtc_signalling"

type rtcSignalling struct {
	log      portsoutlogging.Log
	rtcMedia portsoutperipheral.RTCMedia
}

func New(
	log portsoutlogging.Log,
	rtcMedia portsoutperipheral.RTCMedia,
) portsinhttp.RTCSignalling {
	return &rtcSignalling{
		log:      log,
		rtcMedia: rtcMedia,
	}
}

func (r *rtcSignalling) CreatePeer(ctx context.Context, peerID string, onICECandidate func(domainmodel.ICECandidateData)) (offerSDP string, err error) {
	const tag = path + "/CreatePeer"

	offerSDP, err = r.rtcMedia.CreatePeer(ctx, peerID, onICECandidate)
	if err != nil {
		r.log.Error(
			ctx, tag,
			"Failed to create RTC peer",
			domainmodel.LogMeta{
				"peer_id": peerID,
			},
		)
		return
	}

	r.log.Info(
		ctx, tag,
		"Peer created successfully",
		domainmodel.LogMeta{
			"peer_id": peerID,
		},
	)

	return
}

func (r *rtcSignalling) SetAnswer(ctx context.Context, peerID string, answerSDP string) (err error) {
	const tag = path + "/SetAnswer"

	err = r.rtcMedia.SetAnswer(ctx, peerID, answerSDP)
	if err != nil {
		r.log.Error(
			ctx, tag,
			"Failed to set answer",
			domainmodel.LogMeta{
				"peer_id":    peerID,
				"answer_sdp": answerSDP,
			},
		)
		return
	}

	r.log.Info(
		ctx, tag,
		"Answer set successfully",
		domainmodel.LogMeta{
			"peer_id":    peerID,
			"answer_sdp": answerSDP,
		},
	)

	return
}

func (r *rtcSignalling) AddICECandidate(ctx context.Context, peerID string, candidate domainmodel.ICECandidateData) (err error) {
	const tag = path + "/AddICECandidate"

	err = r.rtcMedia.AddICECandidate(ctx, peerID, candidate)
	if err != nil {
		r.log.Error(
			ctx, tag,
			"Failed to set ICE candidate",
			domainmodel.LogMeta{
				"peer_id":   peerID,
				"candidate": candidate,
			},
		)
		return
	}

	r.log.Info(
		ctx, tag,
		"ICE candidate added successfully",
		domainmodel.LogMeta{
			"peer_id":   peerID,
			"candidate": candidate,
		},
	)

	return
}

func (r *rtcSignalling) ClosePeer(ctx context.Context, peerID string) (err error) {
	const tag = path + "/ClosePeer"

	err = r.rtcMedia.ClosePeer(ctx, peerID)
	if err != nil {
		r.log.Error(
			ctx, tag,
			"Failed to close peer",
			domainmodel.LogMeta{
				"peer_id": peerID,
			},
		)
	}

	r.log.Info(
		ctx, tag,
		"Peer closed successfully",
		domainmodel.LogMeta{
			"peer_id": peerID,
		},
	)

	return
}
