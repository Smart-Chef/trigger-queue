package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

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
		url := os.Getenv("NLP_API") + "/send_message/" + string(data)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Error(err.Error())
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error(err.Error())
		}
		defer res.Body.Close()
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

// Actions for the trigger-queue to execute
var Actions = map[string]Action{
	"sendToNLP":  sendDataHelper("NLP"),
	"changeStep": changeStep,
	"mockAction": mockAction,
}
