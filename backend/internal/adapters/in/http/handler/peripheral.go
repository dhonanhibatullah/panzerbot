package adaptersinhttphandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	adaptersinhttpdto "github.com/dhonanhibatullah/panzerbot/backend/internal/adapters/in/http/dto"
	portsinhttp "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/in/http"
	portsoutlogging "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/logging"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const pathPeripheral = "adapters/in/http/handler/peripheral"

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Peripheral struct {
	log           portsoutlogging.Log
	peripheralSvc portsinhttp.Peripheral
}

func NewPeripheral(
	log portsoutlogging.Log,
	peripheralSvc portsinhttp.Peripheral,
) *Peripheral {
	return &Peripheral{
		log:           log,
		peripheralSvc: peripheralSvc,
	}
}

// @Summary      Peripheral WebSocket
// @Description  Establish a WebSocket connection to control motor speed and servo angle in real-time
// @Tags         peripheral
// @Accept       json
// @Produce      json
// @Success      101  "Switching Protocols"
// @Failure      500  {object}  map[string]string
// @Router       /v1/peripheral/ws [get]
func (p *Peripheral) Ws(c *gin.Context) {
	const tag = pathPeripheral + "/Ws"

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			adaptersinhttpdto.CommonErrorResponse{
				Error: "Failed to handle websocket request",
			},
		)
		return
	}
	defer conn.Close()

	for {
		var req adaptersinhttpdto.PeripheralWsRequest
		if err := conn.ReadJSON(&req); err != nil {
			break
		}

		switch req.Code {
		case "motor":
			var data adaptersinhttpdto.PeripheralWsMotorControlData
			if err := json.Unmarshal(req.Data, &data); err != nil {
				p.log.Warn(
					c.Request.Context(), tag,
					"Bad servo data received",
					nil,
				)
				conn.WriteJSON(adaptersinhttpdto.CommonErrorResponse{
					Error: "Invalid motor data",
				})
				continue
			}
			if err := p.peripheralSvc.SetMotorSpeedScale(
				c.Request.Context(),
				data.Right,
				data.Left,
			); err != nil {
				conn.WriteJSON(adaptersinhttpdto.CommonErrorResponse{
					Error: "An error occured: " + err.Error(),
				})
			}

		case "servo":
			var data adaptersinhttpdto.PeripheralWsServoControlData
			if err := json.Unmarshal(req.Data, &data); err != nil {
				p.log.Warn(
					c.Request.Context(), tag,
					"Bad servo data received",
					nil,
				)
				conn.WriteJSON(adaptersinhttpdto.CommonErrorResponse{
					Error: "Invalid servo data",
				})
				continue
			}
			if err := p.peripheralSvc.SetServoAngle(
				c.Request.Context(),
				data.Pan,
				data.Tilt,
			); err != nil {
				conn.WriteJSON(adaptersinhttpdto.CommonErrorResponse{
					Error: "An error occured: " + err.Error(),
				})
			}

		default:
			p.log.Warn(
				c.Request.Context(), tag,
				"Unknown code detected",
				nil,
			)
			conn.WriteJSON(adaptersinhttpdto.CommonErrorResponse{
				Error: "Unknown code",
			})
		}
	}
}

// @Summary      Get soundboard tracks
// @Description  Returns the list of available soundboard track names
// @Tags         peripheral
// @Produce      json
// @Success      200  {object}  adaptersinhttpdto.PeripheralGetSoundboardTrackResponse
// @Failure      500  {object}  map[string]string
// @Router       /v1/peripheral/soundboard [get]
func (p *Peripheral) GetSoundboardTrack(c *gin.Context) {
	trackNames, err := p.peripheralSvc.GetSoundboardTrackNames(c.Request.Context())
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			adaptersinhttpdto.CommonErrorResponse{
				Error: "An error occured: " + err.Error(),
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		adaptersinhttpdto.PeripheralGetSoundboardTrackResponse{
			Tracks: trackNames,
		},
	)
}

// @Summary      Play a soundboard track
// @Description  Plays a soundboard track by its index
// @Tags         peripheral
// @Param        track_idx  path  int  true  "Track index"  example(0)
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /v1/peripheral/soundboard/{track_idx} [post]
func (p *Peripheral) PostSoundboardTrack(c *gin.Context) {
	idx, err := strconv.Atoi(c.Param("track_idx"))
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			adaptersinhttpdto.CommonErrorResponse{
				Error: "Invalid track index",
			},
		)
		return
	}
	if err := p.peripheralSvc.PlaySoundboardTrack(c.Request.Context(), idx); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			adaptersinhttpdto.CommonErrorResponse{
				Error: "Failed to play soundboard track",
			},
		)
		return
	}
	c.Status(http.StatusNoContent)
}

// @Summary      Stop all soundboard tracks
// @Description  Stops all currently playing soundboard tracks
// @Tags         peripheral
// @Success      204  "No Content"
// @Failure      500  {object}  map[string]string
// @Router       /v1/peripheral/soundboard/stop [post]
func (p *Peripheral) PostSoundboardTrackStop(c *gin.Context) {
	if err := p.peripheralSvc.StopSoundboardTracks(c.Request.Context()); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			adaptersinhttpdto.CommonErrorResponse{
				Error: "Failed to stop soundboard tracks",
			},
		)
		return
	}
	c.Status(http.StatusNoContent)
}
