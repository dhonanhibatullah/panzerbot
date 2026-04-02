package adaptersinhttphandler

import (
	"net/http"
	"sync"

	adaptersinhttpdto "github.com/dhonanhibatullah/panzerbot/backend/internal/adapters/in/http/dto"
	domainmodel "github.com/dhonanhibatullah/panzerbot/backend/internal/domain/model"
	portsinhttp "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/in/http"
	portsoutlogging "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/logging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RTCSignalling struct {
	log              portsoutlogging.Log
	rtcSignallingSvc portsinhttp.RTCSignalling
}

func NewRTCSignalling(
	log portsoutlogging.Log,
	rtcSignallingSvc portsinhttp.RTCSignalling,
) *RTCSignalling {
	return &RTCSignalling{
		rtcSignallingSvc: rtcSignallingSvc,
		log:              log,
	}
}

// @Summary      RTC Signalling WebSocket
// @Description  Establish a WebSocket connection for WebRTC signalling. The robot sends an offer immediately after connecting. The peer replies with an answer and exchanges ICE candidates to establish the media connection. The robot streams camera and microphone; the peer may send audio back to the robot's speaker.
// @Tags         rtc
// @Accept       json
// @Produce      json
// @Success      101  "Switching Protocols"
// @Failure      500  {object}  adaptersinhttpdto.CommonErrorResponse
// @Router       /v1/rtc/ws [get]
func (h *RTCSignalling) Ws(c *gin.Context) {
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			adaptersinhttpdto.CommonErrorResponse{
				Error: "An error occured: " + err.Error(),
			},
		)
		return
	}
	defer conn.Close()

	peerID := uuid.NewString()
	ctx := c.Request.Context()

	var wsMu sync.Mutex
	safeSend := func(msg adaptersinhttpdto.RtcOutbound) {
		wsMu.Lock()
		defer wsMu.Unlock()
		conn.WriteJSON(msg)
	}

	offerSDP, err := h.rtcSignallingSvc.CreatePeer(ctx, peerID, func(cand domainmodel.ICECandidateData) {
		safeSend(adaptersinhttpdto.RtcOutbound{
			Type:          "ice-candidate",
			Candidate:     cand.Candidate,
			SDPMid:        cand.SDPMid,
			SDPMLineIndex: cand.SDPMLineIndex,
		})
	})
	if err != nil {
		safeSend(adaptersinhttpdto.RtcOutbound{
			Type: "error",
		})
		return
	}

	safeSend(adaptersinhttpdto.RtcOutbound{
		Type: "offer",
		SDP:  offerSDP,
	})

loop:
	for {
		var msg adaptersinhttpdto.RtcInbound
		if err := conn.ReadJSON(&msg); err != nil {
			break loop
		}
		switch msg.Type {
		case "answer":
			h.rtcSignallingSvc.SetAnswer(ctx, peerID, msg.SDP)
		case "ice-candidate":
			h.rtcSignallingSvc.AddICECandidate(
				ctx, peerID,
				domainmodel.ICECandidateData{
					Candidate:     msg.Candidate,
					SDPMid:        msg.SDPMid,
					SDPMLineIndex: msg.SDPMLineIndex,
				},
			)
		case "close":
			break loop
		}
	}

	h.rtcSignallingSvc.ClosePeer(ctx, peerID)
}
