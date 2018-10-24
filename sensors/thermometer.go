package sensors

import "sync"

var thermometerInstance *Thermometer
var thermometerOnce sync.Once

// Scale should be treated as a singleton
type Thermometer struct {
	name string
	// TODO: add other metadata
}

// setupThermometer connects to the physical thermometer
func (Thermometer) setupThermometer() *Thermometer {
	// TODO: add code to connect to the actual sensor
	return &Thermometer{"testThermometer"}
}

// GetTemp gets the current temperature value from the thermometer
func (*Thermometer) GetTemp() int {
	// TODO: get temp from sensor
	return 200
}

// Implement Singleton GetInstance
func (*Thermometer) GetInstance() *Thermometer {
	thermometerOnce.Do(func() {
		thermometerInstance = new(Thermometer).setupThermometer()
	})
	return thermometerInstance
}
