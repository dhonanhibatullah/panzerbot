package portsoutperipheral

import (
	"context"

	domainmodel "github.com/dhonanhibatullah/panzerbot/backend/internal/domain/model"
)

type RTCMedia interface {
	CreatePeer(ctx context.Context, peerID string, onICECandidate func(domainmodel.ICECandidateData)) (offerSDP string, err error)
	SetAnswer(ctx context.Context, peerID string, answerSDP string) (err error)
	AddICECandidate(ctx context.Context, peerID string, candidate domainmodel.ICECandidateData) (err error)
	ClosePeer(ctx context.Context, peerID string) (err error)
}
