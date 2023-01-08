// Package device implements different functionality around GPIO attached devices of the raspberry pi
//
// It's based of stianeikeland's go-rpio libary
// See https://github.com/stianeikeland/go-rpio
// It abstracts interaction with the GPIO devices for the main libray,
// so that this package could theoretically use another GPIO libary without chainching the main programm
package device

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	rpio "github.com/stianeikeland/go-rpio/v4"
)

func init() {
	// open memory-range for GPIO access (device /dev/gpiomem)
	// note: user should be member of the gpio group for this to work
	err := rpio.Open()
	if err != nil {
		log.Fatal(err)
	}
}

type Device struct {
	LEDs    map[string]LED
	Buttons map[string]Button
}

type LED struct {
	Color string
	rpio.Pin
}

type Button struct {
	EnableFallEdgeDetection bool
	rpio.Pin
}

// New initializes the given GPIO devices according to their type
func New(leds map[string]LED, buttons map[string]Button) *Device {
	// configure LEDs
	for name, led := range leds {
		led.Output()
		led.PullDown()
		log.Debugf("turned-off LED %s", name)
	}
	// configure buttons
	for name, button := range buttons {
		button.Input()
		log.Debugf("set button %s as input", name)
		if button.EnableFallEdgeDetection {
			button.Detect(rpio.FallEdge) // enable falling edge event detection
			log.Debugf("enabled falling edge detection for '%s'", name)
		}
	}
	log.Print("finished device initialization")
	return &Device{
		LEDs:    leds,
		Buttons: buttons,
	}
}

// Close finalizes all GPIO devices and related memory
func (d *Device) Close() {
	// disable eventdetection on all button when exiting
	for _, button := range d.Buttons {
		if button.EnableFallEdgeDetection {
			button.Detect(rpio.NoEdge)
		}
	}

	// turn-off all LEDs when closing programm
	for _, led := range d.LEDs {
		led.PullDown()
	}

	// Unmap gpio memory when done
	rpio.Close()
}

// Blink starts blinking the given LED until context is closed
func (d *Device) Blink(ctx context.Context, led string, blinkInterval time.Duration) {
	for {
		select {
		case <-ctx.Done():
			d.LEDs[led].PullDown()
			return
		default:
			d.LEDs[led].Toggle()
		}
		time.Sleep(blinkInterval)
	}
}

// BlinkDuration blinks for the given duration
func (d *Device) BlinkDuration(ctx context.Context, led string, blinkInterval, duration time.Duration) {
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			d.LEDs[led].PullDown()
			return
		default:
			d.LEDs[led].Toggle()
			time.Sleep(blinkInterval)
		}
	}
}

// ButtonEvents writes to a channel whenever a button is pressed
func (d *Device) ButtonEvents(ctx context.Context, button string, notify chan bool, sleepDuration time.Duration) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if d.Buttons[button].EdgeDetected() { // check if event occured
				notify <- true
			}
			time.Sleep(sleepDuration)
		}
	}
}
