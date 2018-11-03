package sensors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThermometer_GetInstance(t *testing.T) {
	a := new(Thermometer).GetInstance()
	b := new(Thermometer).GetInstance()

	assert.Equal(t, &a, &b)
}

func TestThermometer_GetTemp(t *testing.T) {
	s := new(Thermometer).GetInstance()
	assert.Equal(t, float64(200), s.GetTemp())
}

func TestSomeInterface_SomeFunc(t *testing.T) {
	s := new(Thermometer).GetInstance()
	assert.Equal(t, float64(200), s.GetTemp())
	assert.Equal(t, 300, 300)
}
