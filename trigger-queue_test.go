package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJob_MarshalJSON(t *testing.T) {
	j := Job{
		123,
		make([]Trigger, 0),
		make([]string, 0),
		make([]interface{}, 0),
		nil,
		"",
		"",
		"TestSubscriber",
		time.Now(),
	}
	actual, _ := j.MarshalJSON()
	expected := "{\"id\":123,\"trigger_keys\":[],\"trigger_params\":[],\"action_key\":\"\",\"action_params\":\"\",\"subscriber\":\"TestSubscriber\",\"created_at\":\"2018-10-31T15:42:56.473957-05:00\"}"
	//assert.Equal(t, expected, actual)
	//assert.JSONEq(t, expected, string(actual))
	assert.JSONEqf(t, expected, string(actual), "test")
	//assert.Implementsf()
}

func TestCreateRoutes(t *testing.T) {

}

func TestRouteWalker(t *testing.T) {

}
