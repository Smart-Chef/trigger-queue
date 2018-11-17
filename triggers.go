package main

import (
	"time"
	"trigger-queue/sensors"

	log "github.com/sirupsen/logrus"
)

type Trigger func(interface{}) (bool, error)

func compareValHelper(getVal func() (float64, error), compare func(float64, float64) bool) Trigger {
	return func(val interface{}) (bool, error) {
		sensorVal, err := getVal()
		if err != nil {
			log.Error("Error getting value from sensor")
			log.Error(err)
			return false, err
		}
		return compare(sensorVal, val.(float64)), nil
	}

}

func compareSensorReading(t string, getVal func() (float64, error)) Trigger {
	// May need to modify for IEEE 737 float issues
	// https://floating-point-gui.de/errors/comparison/#look-out-for-edge-cases
	switch t {
	case ">":
		return compareValHelper(getVal, func(f float64, f2 float64) bool {
			return f > f2
		})
	case ">=":
		return compareValHelper(getVal, func(f float64, f2 float64) bool {
			return f >= f2
		})
	case "<":
		return compareValHelper(getVal, func(f float64, f2 float64) bool {
			return f < f2
		})
	case "<=":
		return compareValHelper(getVal, func(f float64, f2 float64) bool {
			return f <= f2
		})
	case "==":
		return compareValHelper(getVal, func(f float64, f2 float64) bool {
			return f == f2
		})
	default:
		return func(val interface{}) (bool, error) {
			return true, nil
		}
	}
}

func timer(val interface{}) (bool, error) {
	t, err := time.Parse(time.RFC3339, val.(string))
	if err != nil {
		log.Error(err)
		return false, err
	}
	loc, _ := time.LoadLocation("UTC")
	return time.Now().In(loc).After(t), nil
}

func tempComparisonHelper(t string, thermometer *sensors.Thermometer) Trigger {
	return compareSensorReading(t, thermometer.GetTemp)
}

func weightComparisonHelper(t string, scale *sensors.Scale) Trigger {
	return compareSensorReading(t, scale.GetWeight)
}

var Triggers = map[string]Trigger{
	"timer":     timer,
	"temp_>":    tempComparisonHelper(">", Thermometer),
	"temp_>=":   tempComparisonHelper(">=", Thermometer),
	"temp_<":    tempComparisonHelper("<", Thermometer),
	"temp_<=":   tempComparisonHelper("<=", Thermometer),
	"temp_==":   tempComparisonHelper("==", Thermometer),
	"weight_>":  weightComparisonHelper(">", Scale),
	"weight_>=": weightComparisonHelper(">=", Scale),
	"weight_<":  weightComparisonHelper("<", Scale),
	"weight_<=": weightComparisonHelper("<=", Scale),
	"weight_==": weightComparisonHelper("==", Scale),
}
