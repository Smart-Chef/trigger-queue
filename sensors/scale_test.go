package sensors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScale_GetInstance(t *testing.T) {
	a := new(Scale).GetInstance()
	b := new(Scale).GetInstance()

	assert.Equal(t, &a, &b)
}
