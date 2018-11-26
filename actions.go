package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Action func(interface{})

func sendDataHelper(service string) Action {
	return func(d interface{}) {
		data, err := json.Marshal(d)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("Sending the following to " + service)
		log.Info(string(data))
		// TODO: actually send the data
	}
}

// Send payload to recipe-walkthrough endpoint
func changeStep(payload interface{}) {
	url := os.Getenv("RECIPE_WALKTHROUGH_API") + "/increment-step"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload.(string))))
	if err != nil {
		log.Error("Could not execute action changeStep")
		log.Error(err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("Could not execute action changeStep")
		log.Error(err.Error())
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("Error reading response")
		log.Error(err.Error())
	} else {
		log.Info("Action response: " + string(body))
	}
}

// Action used for testing purposes that just logs the payload
func mockAction(i interface{}) {
	log.Info(i.(string))
}

func setStoveTemp(i interface{}) {
	Stove.SetTemp(int(i.(float64)))
}

func changeStepStove(i interface{}) {
	type Data struct {
		IncrementSteps string `json:"increment_steps"`
		SendToNLP      bool   `json:"send_to_nlp"`
		StoveStart     bool   `json:"stove_start"`
		StoveTemp      string `json:"stove_temp"`
	}

	// Check what stove changes needs to be done
	buffer := []byte(i.(string))
	var data Data
	e := json.Unmarshal(buffer, &data)
	if e != nil {
		log.Error(e.Error())
	}

	temp, e := strconv.ParseFloat(data.StoveTemp, 64)
	if e != nil {
		log.Error(e.Error())
	}
	if e = Stove.SetTemp(int(temp)); e != nil {
		log.Error(e.Error())
	}

	if data.StoveStart {
		if e = Stove.StartStove(); e != nil {
			log.Error(e.Error())
		}
	}
	// Change the step like normal
	changeStep(i)
}

// Actions for the trigger-queue to execute
var Actions = map[string]Action{
	"setStoveTemp":    setStoveTemp,
	"sendToNLP":       sendDataHelper("NLP"),
	"changeStep":      changeStep,
	"changeStepStove": changeStepStove,
	"mockAction":      mockAction,
}
