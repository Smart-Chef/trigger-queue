package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemp(t *testing.T) {
	tests := []struct {
		name     string
		param    int
		expected bool
	}{
		{">", 200, true},
		{">=", 200, true},
		{"<", 200, false},
		{"<=", 200, false},
		{"==", 200, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, compareSensorReading(test.name, func() float64 {
				return 250
			})(float64(test.param)), test.expected)
		})
	}
}

func TestWeight(t *testing.T) {
	tests := []struct {
		name     string
		param    int
		expected bool
	}{
		{"weight_>", 200, true},
		{"weight_>=", 200, true},
		{"weight_<", 200, false},
		{"weight_<=", 200, false},
		{"weight_==", 200, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, Triggers[test.name](float64(test.param)), test.expected)
		})
	}
}
