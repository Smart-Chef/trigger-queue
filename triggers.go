package main

import (
	"trigger-queue/sensors"
)

type Trigger func(interface{}) bool

var Scale = new(sensors.Scale).GetInstance()
var Thermometer = new(sensors.Thermometer).GetInstance()

func compareSensorReading(t string, getVal func() float64) Trigger {
	switch t {
	case ">":
		return func(val interface{}) bool {
			return val.(float64) > getVal()
		}
	case ">=":
		return func(val interface{}) bool {
			return val.(float64) >= getVal()
		}
	case "<":
		return func(val interface{}) bool {
			return val.(float64) < getVal()
		}
	case "<=":
		return func(val interface{}) bool {
			return val.(float64) <= getVal()
		}
	case "==":
		return func(val interface{}) bool {
			return val.(float64) == getVal()
		}
	default:
		return func(val interface{}) bool {
			return true
		}
	}
}

func tempComparisonHelper(t string, thermometer *sensors.Thermometer) Trigger {
	return compareSensorReading(t, thermometer.GetTemp)
}

func weightComparisonHelper(t string, scale *sensors.Scale) Trigger {
	return compareSensorReading(t, scale.GetWeight)
}

var Triggers = map[string]Trigger{
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
