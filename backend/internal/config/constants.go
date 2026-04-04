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

	ServoPanPin    rpio.Pin = 18
	ServoTiltPin   rpio.Pin = 19
	PwmCycleLength uint32   = 20_000
	PwmFrequency   int      = 50 * int(PwmCycleLength)

	// Linux sysfs PWM chip paths.
	// On modern RPi OS (kernel 5.x+) each PWM channel is its own pwmchip.
	// Verify with: ls /sys/class/pwm/
	ServoPanPwmChip  string = "/sys/class/pwm/pwmchip0" // GPIO 18 (PWM0)
	ServoTiltPwmChip string = "/sys/class/pwm/pwmchip2" // GPIO 19 (PWM1)

	RTCStunServer            string          = "stun:stun.l.google.com:19302"
	RTCVideoBitRate          int             = 1_000_000
	RTCOpusSampleRate        beep.SampleRate = 48000
	RTCOpusSamplesPerChannel int             = 960
	RTCOpusFrameBuffer       int             = 1920
	RTCBeepStreamerBuffer    int             = 4096
)
