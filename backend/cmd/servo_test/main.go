package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Sysfs PWM on Raspberry Pi 4
// Requires in /boot/firmware/config.txt (or /boot/config.txt):
//   dtoverlay=pwm,pin=18,func=2      (single channel, GPIO 18)
//   dtoverlay=pwm-2chan               (GPIO 18 + GPIO 19, channels 0 and 1)
//
// Then: /sys/class/pwm/pwmchip0/pwm0 becomes available.

const (
	pwmChip    = "/sys/class/pwm/pwmchip0"
	pwmChannel = 0 // channel 0 = GPIO 18 (Alt5) or GPIO 12 (Alt0)

	periodNs    = 20_000_000 // 20 ms → 50 Hz
	pulseMinNs  = 500_000    // 500 µs → 0°
	pulseMaxNs  = 2_500_000  // 2500 µs → 180°
	nsPerDeg    = (pulseMaxNs - pulseMinNs) / 180 // ~11,111 ns/degree
)

func pwmPath(channel int) string {
	return fmt.Sprintf("%s/pwm%d", pwmChip, channel)
}

func write(path, value string) error {
	return os.WriteFile(path, []byte(value), 0644)
}

func exportChannel(channel int) error {
	path := pwmPath(channel)
	if _, err := os.Stat(path); err == nil {
		return nil // already exported
	}
	return write(pwmChip+"/export", strconv.Itoa(channel))
}

func enablePwm(channel int, enable bool) error {
	val := "0"
	if enable {
		val = "1"
	}
	return write(pwmPath(channel)+"/enable", val)
}

func setPeriod(channel int, ns int) error {
	return write(pwmPath(channel)+"/period", strconv.Itoa(ns))
}

func setDutyCycle(channel int, ns int) error {
	return write(pwmPath(channel)+"/duty_cycle", strconv.Itoa(ns))
}

func setAngleDeg(channel, deg int) {
	if deg < 0 {
		deg = 0
	}
	if deg > 180 {
		deg = 180
	}
	pulseNs := pulseMinNs + deg*nsPerDeg
	if err := setDutyCycle(channel, pulseNs); err != nil {
		fmt.Fprintf(os.Stderr, "duty_cycle error: %v\n", err)
	}
}

func mustInit(channel int) {
	if err := exportChannel(channel); err != nil {
		fmt.Fprintf(os.Stderr, "export error: %v\n", err)
		os.Exit(1)
	}
	// Disable before changing period (required by kernel driver)
	_ = enablePwm(channel, false)
	if err := setPeriod(channel, periodNs); err != nil {
		fmt.Fprintf(os.Stderr, "period error: %v\n", err)
		os.Exit(1)
	}
	if err := setDutyCycle(channel, pulseMinNs); err != nil {
		fmt.Fprintf(os.Stderr, "duty_cycle error: %v\n", err)
		os.Exit(1)
	}
	if err := enablePwm(channel, true); err != nil {
		fmt.Fprintf(os.Stderr, "enable error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	mustInit(pwmChannel)
	defer enablePwm(pwmChannel, false)

	fmt.Printf("PWM sysfs: channel=%d, period=%dns (~50Hz), step=%dns/deg\n",
		pwmChannel, periodNs, nsPerDeg)

	for {
		fmt.Println("--- sweep forward ---")
		for deg := 0; deg <= 180; deg++ {
			setAngleDeg(pwmChannel, deg)
			fmt.Printf("%3d°\n", deg)
			time.Sleep(50 * time.Millisecond)
		}

		fmt.Println("--- sweep back ---")
		for deg := 180; deg >= 0; deg-- {
			setAngleDeg(pwmChannel, deg)
			fmt.Printf("%3d°\n", deg)
			time.Sleep(50 * time.Millisecond)
		}
	}
}
