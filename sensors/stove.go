package sensors

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio"
)

var (
	stoveInstance *Stove
	stoveOnce     sync.Once

	KeepWarm = 150
	Medium   = 275
	High     = 425
)

var gpioPins = map[string]rpio.Pin{
	"decrease_temp": rpio.Pin(23), // D
	"increase_temp": rpio.Pin(16), // I
	"keep_warm":     rpio.Pin(25), // K
	"medium":        rpio.Pin(21), // M
	"high":          rpio.Pin(18), // H
	"start":         rpio.Pin(12), // S
}

// Stove should be treated as a singleton
type Stove struct {
	name  string
	pins  map[string]rpio.Pin
	setup bool
}

// setupStove connects to the physical sensor
func (Stove) setupStove() (*Stove, error) {
	// Setup rpio
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
	}

	// Set all the pins to High
	for _, pin := range gpioPins {
		pin.Output()
		pin.High()
	}

	log.Info("Done setting up stove")
	return &Stove{
		name:  "stove",
		pins:  gpioPins,
		setup: true,
	}, nil
}

func (s *Stove) StartStove() error {
	if !s.setup {
		return errors.New("This stove instance has not been setup")
	}
	return toggleButton("start")
}

// SetTemp value
func (s *Stove) SetTemp(temp int) error {
	if !s.setup {
		return errors.New("This stove instance has not been setup")
	}

	if temp <= KeepWarm {
		toggleButton("keep_warm")
		decTimes := math.Ceil(float64(KeepWarm-temp) / 10.0)

		// Decreace to a minimum 100 degrees
		for i := 0; i < int(decTimes) && i < 5; i++ {
			toggleButton("decrease_temp")
		}
	} else if temp <= Medium {
		if (Medium - temp) > (temp - KeepWarm) {
			toggleButton("keep_warm")
			incTimes := math.Ceil(float64(temp-KeepWarm) / 10)
			for i := 0; i < int(incTimes); i++ {
				toggleButton("increase_temp")
			}
		} else {
			toggleButton("medium")
			decTimes := math.Ceil(float64(Medium-temp) / 10)
			for i := 0; i < int(decTimes); i++ {
				toggleButton("decrease_temp")
			}
		}
	} else if temp <= High {
		if (High - temp) > (temp - Medium) {
			toggleButton("medium")

			incTimes := math.Ceil(float64(temp-Medium) / 10)
			for i := 0; i < int(incTimes); i++ {
				toggleButton("increase_temp")
			}
		} else {
			toggleButton("high")
			decTimes := math.Ceil(float64(High-temp) / 10)
			for i := 0; i < int(decTimes); i++ {
				toggleButton("decrease_temp")
			}
		}
	} else if temp > High {
		toggleButton("high")
		incTimes := math.Ceil(float64(temp-High) / 10)
		// Break at maximum 500 degrees
		for i := 0; i < int(incTimes) && i < 8; i++ {
			toggleButton("increase_temp")
		}
	}
	return nil
}

// Cleanup the rpio
func (*Stove) Cleanup() error {
	return rpio.Close()
}

// Implement Singleton GetInstance
func (*Stove) GetInstance() (*Stove, error) {
	var err error
	err = nil
	stoveOnce.Do(func() {
		stoveInstance, err = new(Stove).setupStove()
	})
	return stoveInstance, err
}

func (s *Stove) toggleAll() error {
	pins := []rpio.Pin{rpio.Pin(12), rpio.Pin(16), rpio.Pin(18), rpio.Pin(21), rpio.Pin(23), rpio.Pin(25)}
	for _, pin := range pins {
		fmt.Println(pin)
		pin.Low()
		time.Sleep(time.Second / 10)
		pin.High()
		time.Sleep(time.Second)
	}
	return nil
}

// toggleButton - press and depressed btn
func toggleButton(key string) error {
	btn, ok := gpioPins[key]
	if !ok {
		return errors.New("button not found: " + key)
	}
	fmt.Printf("Toggling: %s - %d\n", key, btn)
	btn.Low()
	time.Sleep(time.Second / 10)
	btn.High()
	time.Sleep(time.Second)
	return nil
}
