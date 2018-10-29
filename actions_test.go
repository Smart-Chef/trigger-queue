package main

import "testing"

func Test_sendData(t *testing.T) {
	type Payload struct {
		Status string `json:"status"`
		Data   string `json:"data"`
	}
	Actions["sendToNLP"](Payload{"success", "200"})
}
