package adaptersinhttpdto

import "encoding/json"

type PeripheralWsRequest struct {
	Code string          `json:"code" binding:"required" example:"motor"`
	Data json.RawMessage `json:"data" binding:"required" swaggertype:"object" example:"{\"right\":0.7,\"left\":-1.0}"`
}

type PeripheralWsMotorControlData struct {
	Right float64 `json:"right" binding:"required" example:"0.7"`
	Left  float64 `json:"left" binding:"required" example:"-1.0"`
}

type PeripheralWsServoControlData struct {
	Pan  float64 `json:"pan" binding:"required" example:"1.292"`
	Tilt float64 `json:"tilt" binding:"required" example:"0.35"`
}

type PeripheralGetSoundboardTrackResponse struct {
	Tracks []string `json:"tracks" example:"Boom,Fart,Get Out,Oi Oi Oi,Outro,Rizz,Uwu"`
}
