package sensors

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/stianeikeland/go-rpio"
)

var (
	scaleInstance *Scale
	scaleOnce     sync.Once

	KeepWarm = 150
	Medium   = 275
	High     = 425
)

var gpioPins = map[string]rpio.Pin{
	"keep_warm":     rpio.Pin(18),
	"medium":        rpio.Pin(18),
	"high":          rpio.Pin(18),
	"increase_temp": rpio.Pin(18),
	"decrease_temp": rpio.Pin(18),
	"start":         rpio.Pin(18),
}

// Scale should be treated as a singleton
type Scale struct {
	name  string
	pins  map[string]rpio.Pin
	setup bool
}

// setupScale connects to the physical sensor
func (Scale) setupScale() (*Scale, error) {
	// Setup rpio
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
	}

	// Set all the pins to High
	for _, pin := range gpioPins {
		pin.High()
	}

	fmt.Println("Done with setup")
	return &Scale{
		name:  "scale",
		pins:  gpioPins,
		setup: true,
	}, nil
}

// SetTemp value
func (s *Scale) SetTemp(temp int) error {
	if !s.setup {
		return errors.New("This scale instance has not been setup")
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
		for i := 0; i < int(incTimes) && i < 8; i++ {
			toggleButton("increase_temp")
		}
	}
	return nil
}

// Cleanup the rpio
func (*Scale) Cleanup() error {
	return rpio.Close()
}

// GetWeight gets the current weight value from teh scale sensor
func (*Scale) GetWeight() (float64, error) {
	// TODO: get weight from sensor
	return 250, nil
}

// Implement Singleton GetInstance
func (*Scale) GetInstance() (*Scale, error) {
	var err error
	err = nil
	scaleOnce.Do(func() {
		scaleInstance, err = new(Scale).setupScale()
	})
	return scaleInstance, err
}

// toggleButton - press and depressed btn
func toggleButton(key string) error {
	fmt.Println("Looking for key")
	btn, ok := gpioPins[key]
	if !ok {
		return errors.New("button not found: " + key)
	}

	fmt.Println("Button low")
	btn.Low()
	time.Sleep(time.Second / 10)
	btn.High()
	fmt.Println("Button High")
	time.Sleep(time.Second)
	return nil
}
