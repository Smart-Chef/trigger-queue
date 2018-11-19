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
func (Scale) setupScale() (*Scale, error) {
	// TODO: add code to connect to the actual sensor
	return &Scale{"testScale"}, nil
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
