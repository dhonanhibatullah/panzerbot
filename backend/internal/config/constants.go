package config

import (
	"log/slog"

	"github.com/gopxl/beep"
	"github.com/peergum/go-rpio/v5"
)

var (
	LogLevel slog.Level = slog.LevelInfo

	// MotorRightAPin rpio.Pin = 7
	// MotorRightBPin rpio.Pin = 1
	// MotorRightPwmPin rpio.Pin = 12
	// MotorLeftAPin rpio.Pin = 5
	// MotorLeftBPin rpio.Pin = 6
	// MotorLeftPwmPin  rpio.Pin = 13
	// PwmFrequency   int    = 64000
	// PwmCycleLength uint32 = 1024

	// ServoPanPin  rpio.Pin = 27
	// ServoTiltPin rpio.Pin = 22

	MotorRightAPin rpio.Pin = 24
	MotorRightBPin rpio.Pin = 23
	MotorLeftAPin  rpio.Pin = 22
	MotorLeftBPin  rpio.Pin = 27

	ServoPanPin    rpio.Pin = 12
	ServoTiltPin   rpio.Pin = 13
	PwmCycleLength uint32   = 20_000
	PwmFrequency   int      = 50 * int(PwmCycleLength)

	RTCStunServer            string          = "stun:stun.l.google.com:19302"
	RTCVideoBitRate          int             = 1_000_000
	RTCOpusSampleRate        beep.SampleRate = 48000
	RTCOpusSamplesPerChannel int             = 960
	RTCOpusFrameBuffer       int             = 1920
	RTCBeepStreamerBuffer    int             = 4096
)
