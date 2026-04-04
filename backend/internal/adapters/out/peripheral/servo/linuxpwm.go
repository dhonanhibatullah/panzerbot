package adaptersoutperipheralservo

import (
	"context"
	"fmt"
	"math"
	"os"
	"strconv"

	domainmodel "github.com/dhonanhibatullah/panzerbot/backend/internal/domain/model"
	portsoutlogging "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/logging"
	portsoutperipheral "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/peripheral"
)

const (
	periodNs   = 20_000_000 // 20 ms → 50 Hz
	pulseMinNs = 500_000    // 500 µs → 0°
	pulseMaxNs = 2_500_000  // 2500 µs → 180°
)

type servoLinuxPwm struct {
	log      portsoutlogging.Log
	chipPath string // e.g. /sys/class/pwm/pwmchip0 or /sys/class/pwm/pwmchip2
	channel  int    // on modern kernels (5.x+) each chip has only channel 0
}

// NewLinuxPwm controls a servo via the Linux sysfs PWM interface.
//
// On modern Raspberry Pi OS (kernel 5.x+), each PWM channel gets its own
// pwmchip device instead of sharing one:
//
//	GPIO 18 (PWM0) → chipPath="/sys/class/pwm/pwmchip0", channel=0
//	GPIO 19 (PWM1) → chipPath="/sys/class/pwm/pwmchip2", channel=0
//
// Verify on the Pi with: ls /sys/class/pwm/
//
// Requires in /boot/firmware/config.txt:
//
//	dtoverlay=pwm-2chan,pin=18,func=2,pin2=19,func2=2

func NewLinuxPwm(
	log portsoutlogging.Log,
	chipPath string,
	channel int,
) portsoutperipheral.Servo {
	s := &servoLinuxPwm{
		log:      log,
		chipPath: chipPath,
		channel:  channel,
	}
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
		s.log.Error(
			ctx, tag,
			"Failed to set duty cycle",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
		return err
	}
	return nil
}

func (s *servoLinuxPwm) channelPath() string {
	return fmt.Sprintf("%s/pwm%d", s.chipPath, s.channel)
}

func (s *servoLinuxPwm) writeFile(name, value string) error {
	return os.WriteFile(s.channelPath()+"/"+name, []byte(value), 0644)
}

func (s *servoLinuxPwm) mustInit() {
	if _, err := os.Stat(s.channelPath()); err != nil {
		if err := os.WriteFile(s.chipPath+"/export", []byte(strconv.Itoa(s.channel)), 0644); err != nil {
			panic(fmt.Sprintf("servo linuxpwm: export channel %d on %s: %v", s.channel, s.chipPath, err))
		}
	}

	_ = s.writeFile("enable", "0")

	if err := s.writeFile("period", strconv.Itoa(periodNs)); err != nil {
		panic(fmt.Sprintf("servo linuxpwm: set period channel %d on %s: %v", s.channel, s.chipPath, err))
	}
	if err := s.writeFile("duty_cycle", strconv.Itoa(pulseMinNs)); err != nil {
		panic(fmt.Sprintf("servo linuxpwm: set duty_cycle channel %d on %s: %v", s.channel, s.chipPath, err))
	}
	if err := s.writeFile("enable", "1"); err != nil {
		panic(fmt.Sprintf("servo linuxpwm: enable channel %d on %s: %v", s.channel, s.chipPath, err))
	}
}
