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
