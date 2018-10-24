package sensors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThermometer_GetInstance(t *testing.T) {
	a := new(Scale).GetInstance()
	b := new(Scale).GetInstance()

	assert.Equal(t, &a, &b)
}

func TestThermometer_GetTemp(t *testing.T) {
	s := new(Scale).GetInstance()
	assert.Equal(t, s.GetWeight(), 200)
}
