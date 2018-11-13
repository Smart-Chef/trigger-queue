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
		// TODO: actually send the data
	}
}

func changeStep(payload interface{}) {
	url := os.Getenv("RECIPE_WALKTHROUGH_API") + "/increment-step"

	log.Info("Executing \"changeStep\"")
	jstStr, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jstStr))
	req.Header.Add("Content-Type", "application/json")
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	log.Info("Action response: " + string(body))
}

var Actions = map[string]Action{
	"sendToNLP":  sendDataHelper("NLP"),
	"changeStep": changeStep,
}
