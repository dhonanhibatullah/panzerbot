package adaptersoutperipheralservo

import (
	"context"
	"fmt"
	"math"
	"os"
	"strconv"

	portsoutlogging "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/logging"
	portsoutperipheral "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/peripheral"
)

const (
	pwmChipPath = "/sys/class/pwm/pwmchip0"
	periodNs    = 20_000_000 // 20 ms → 50 Hz
	pulseMinNs  = 500_000    // 500 µs → 0°
	pulseMaxNs  = 2_500_000  // 2500 µs → 180°
)

type servoLinuxPwm struct {
	log     portsoutlogging.Log
	channel int
}

func NewLinuxPwm(
	log portsoutlogging.Log,
	channel int,
) portsoutperipheral.Servo {
	s := &servoLinuxPwm{
		log:     log,
		channel: channel}
	s.mustInit()
	return s
}

func (s *servoLinuxPwm) SetAngle(ctx context.Context, angle float64) (err error) {
	const tag = path + "/SetAngle"

	if angle < 0 {
		angle = 0
	} else if angle > math.Pi {
		angle = math.Pi
	}

	pulseNs := int(pulseMinNs + (angle/math.Pi)*float64(pulseMaxNs-pulseMinNs))
	if err = s.writeFile("duty_cycle", strconv.Itoa(pulseNs)); err != nil {
		s.log.Error(ctx, tag, "Failed to set duty cycle", nil)
		return err
	}
	return nil
}

func (s *servoLinuxPwm) channelPath() string {
	return fmt.Sprintf("%s/pwm%d", pwmChipPath, s.channel)
}

func (s *servoLinuxPwm) writeFile(name, value string) error {
	return os.WriteFile(s.channelPath()+"/"+name, []byte(value), 0644)
}

func (s *servoLinuxPwm) mustInit() {
	if _, err := os.Stat(s.channelPath()); err != nil {
		if err := os.WriteFile(pwmChipPath+"/export", []byte(strconv.Itoa(s.channel)), 0644); err != nil {
			panic(fmt.Sprintf("servo linuxpwm: export channel %d: %v", s.channel, err))
		}
	}

	_ = s.writeFile("enable", "0")

	if err := s.writeFile("period", strconv.Itoa(periodNs)); err != nil {
		panic(fmt.Sprintf("servo linuxpwm: set period channel %d: %v", s.channel, err))
	}
	if err := s.writeFile("duty_cycle", strconv.Itoa(pulseMinNs)); err != nil {
		panic(fmt.Sprintf("servo linuxpwm: set duty_cycle channel %d: %v", s.channel, err))
	}
	if err := s.writeFile("enable", "1"); err != nil {
		panic(fmt.Sprintf("servo linuxpwm: enable channel %d: %v", s.channel, err))
	}
}
