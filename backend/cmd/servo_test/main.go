package main

import (
	"fmt"
	"os"
	"time"

	"github.com/peergum/go-rpio/v5"
)

func main() {
	if err := rpio.Open(); err != nil {
		os.Exit(1)
	}
	defer rpio.Close()

	// Pin 18 supports Hardware PWM
	pin := rpio.Pin(18)
	pin.Mode(rpio.Pwm)

	// Set Frequency to 50Hz
	// Formula: 19.2MHz / 64 / 6000 = 50Hz
	pin.Freq(1_000_000)

	// Move to 0 degrees (0.5ms)
	pin.DutyCycle(500, 20_000)
	fmt.Print("0 DEGREE")
	time.Sleep(3 * time.Second)

	// Move to 90 degrees (1.5ms)
	pin.DutyCycle(1500, 20_000)
	fmt.Print("90 DEGREE")
	time.Sleep(3 * time.Second)

	// Move to 180 degrees (2.5ms)
	pin.DutyCycle(2500, 20_000)
	fmt.Print("180 DEGREE")
	time.Sleep(3 * time.Second)
}
