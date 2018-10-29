package main

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Action func(interface{})

func sendDataHelper(service string) Action {
	return func(d interface{}) {
		data, err := json.Marshal(d)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Sending the following to " + service)
		fmt.Printf("%s\n", data)
	}
}

var Actions = map[string]Action{
	"sendToNLP": sendDataHelper("NLP"),
}
