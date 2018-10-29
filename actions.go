package main

import (
	"encoding/json"

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
		log.Info("%s\n", string(data))
		// TODO: actually send the data
	}
}

var Actions = map[string]Action{
	"sendToNLP": sendDataHelper("NLP"),
}
