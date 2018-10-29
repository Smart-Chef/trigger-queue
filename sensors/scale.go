package sensors

import (
	"sync"
)

var scaleInstance *Scale
var scaleOnce sync.Once

// Scale should be treated as a singleton
type Scale struct {
	name string
	// TODO: add other metadata
}

// setupScale connects to the physical sensor
func (Scale) setupScale() *Scale {
	// TODO: add code to connect to the actual sensor
	return &Scale{"testScale"}
}

// GetWeight gets the current weight value from teh scale sensor
func (*Scale) GetWeight() float64 {
	// TODO: get weight from sensor
	return 200
}

// Implement Singleton GetInstance
func (*Scale) GetInstance() *Scale {
	scaleOnce.Do(func() {
		scaleInstance = new(Scale).setupScale()
	})
	return scaleInstance
}
