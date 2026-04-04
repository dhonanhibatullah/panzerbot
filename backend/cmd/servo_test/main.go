package main

import (
	"fmt"
	"os"
	"time"

	"github.com/peergum/go-rpio/v5"
)

const (
	pwmFreq      = 1_000_000                           // 1 MHz clock → 1 µs per tick
	pwmCycleLen  = 20_000                              // 20 ms period → 50 Hz
	pulseMinTick = 500                                 // 500 µs → 0°
	pulseMaxTick = 2500                                // 2500 µs → 180°
	tickPerDeg   = (pulseMaxTick - pulseMinTick) / 180 // ~11 ticks/degree
)

func tickForDeg(deg int) uint32 {
	if deg < 0 {
		deg = 0
	}
	if deg > 180 {
		deg = 180
	}
	return uint32(pulseMinTick + deg*tickPerDeg)
}

func main() {
	if err := rpio.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "rpio.Open failed (try sudo):", err)
		os.Exit(1)
	}
	defer rpio.Close()

	pin := rpio.Pin(18)
	pin.Pwm()
	pin.Freq(pwmFreq)

	fmt.Printf("PWM: freq=%d Hz, cycleLen=%d, tickPerDeg=%d\n",
		pwmFreq, pwmCycleLen, tickPerDeg)

	// Sweep 0° → 180° → 0° in 1° steps
	for {
		fmt.Println("--- sweep forward ---")
		for deg := 0; deg <= 180; deg++ {
			tick := tickForDeg(deg)
			pin.DutyCycle(tick, pwmCycleLen)
			fmt.Printf("%3d° → tick %d\n", deg, tick)
			time.Sleep(50 * time.Millisecond)
		}

		fmt.Println("--- sweep back ---")
		for deg := 180; deg >= 0; deg-- {
			tick := tickForDeg(deg)
			pin.DutyCycle(tick, pwmCycleLen)
			fmt.Printf("%3d° → tick %d\n", deg, tick)
			time.Sleep(50 * time.Millisecond)
		}
	}
}
