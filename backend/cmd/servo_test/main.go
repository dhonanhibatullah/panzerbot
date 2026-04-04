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

	pin := rpio.Pin(18)
	pin.Pwm()
	pin.Freq(1_000_000)

	for {
		pin.DutyCycle(500, 20_000)
		fmt.Println("0 DEGREE")
		time.Sleep(2 * time.Second)

		pin.DutyCycle(1500, 20_000)
		fmt.Println("90 DEGREE")
		time.Sleep(2 * time.Second)

		pin.DutyCycle(2500, 20_000)
		fmt.Println("180 DEGREE")
		time.Sleep(2 * time.Second)
	}
}
