package main

import (
	"time"
	"trigger-queue/sensors"

	log "github.com/sirupsen/logrus"
)

type Trigger func(interface{}) bool

func compareSensorReading(t string, getVal func() float64) Trigger {
	switch t {
	case ">":
		return func(val interface{}) bool {
			return getVal() > val.(float64)
		}
	case ">=":
		return func(val interface{}) bool {
			return getVal() >= val.(float64)
		}
	case "<":
		return func(val interface{}) bool {
			return getVal() < val.(float64)
		}
	case "<=":
		return func(val interface{}) bool {
			return getVal() <= val.(float64)
		}
	case "==":
		return func(val interface{}) bool {
			return getVal() == val.(float64)
		}
	default:
		return func(val interface{}) bool {
			return true
		}
	}
}

func timer(val interface{}) bool {
	t, err := time.Parse(time.RFC3339, val.(string))
	loc, _ := time.LoadLocation("UTC")

	if err != nil {
		log.Error(err.Error())
		return false
	}
	return time.Now().In(loc).After(t)
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
