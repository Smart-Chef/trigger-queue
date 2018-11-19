package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"trigger-queue/queue"
	"trigger-queue/sensors"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var InvalidStatus = "INVALID"

// handleBadRequest for the application
func handleBadRequest(w http.ResponseWriter, msgAndArgs ...interface{}) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(&struct {
		Status string      `json:"status"`
		Msg    interface{} `json:"msg"`
	}{InvalidStatus, msgAndArgs})
}

// addJob to the requested service
var addJob = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type Request struct {
		Service       string        `json:"service"`
		TriggerKeys   []string      `json:"trigger_keys"`
		TriggerParams []interface{} `json:"trigger_params"`
		ActionKey     string        `json:"action_key"`
		ActionParams  interface{}   `json:"action_params"`
	}
	type Response struct {
		Status string `json:"status"`
		JobID  int64  `json:"job_id"`
		Msg    string `json:"msg"`
	}

	req := &Request{}
	json.NewDecoder(r.Body).Decode(&req)

	if req.Service == "" {
		req.Service = "other"
	}

	q, ok := TriggerQueue[req.Service]
	if !ok {
		log.Error("Unrecognized Service \"" + req.Service + "\"")
		handleBadRequest(w, "Unrecognized Service \""+req.Service+"\"")
		return
	}

	j, err := createJob(req.Service, req.TriggerKeys, req.TriggerParams, req.ActionKey, req.ActionParams)
	if err != nil {
		log.Error("Error creating job")
		log.Error(err.Error())
		handleBadRequest(w, "Error creating job")
		return
	}

	id := q.Append(j)
	j.ID = id
	json.NewEncoder(w).Encode(&Response{
		Status: req.Service,
		JobID:  j.ID,
	})
})

// Get the current thermometer value
var GetTemp = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type Response struct {
		Status string  `json:"status"`
		Data   float64 `json:"data"`
	}

	temp, err := new(sensors.Thermometer).GetInstance()
	if err != nil {
		log.Error("Error getting thermometer instance")
		log.Error(err.Error())
		handleBadRequest(w, "Error getting thermometer instance")
	}

	val, err := temp.GetTemp()
	if err != nil {
		log.Error("Error getting thermometer value")
		log.Error(err.Error())
		handleBadRequest(w, "Error getting thermometer value")
	}
	json.NewEncoder(w).Encode(&Response{
		Status: "success",
		Data:   val,
	})
})

// Get the current scale reading
var GetWeight = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type Response struct {
		Status string  `json:"status"`
		Data   float64 `json:"data"`
	}

	scale, err := new(sensors.Scale).GetInstance()
	if err != nil {
		log.Error("Error getting scale instance")
		log.Error(err.Error())
		handleBadRequest(w, "Error getting scale instance")
	}

	val, err := scale.GetWeight()
	if err != nil {
		log.Error("Error getting scale value")
		log.Error(err.Error())
		handleBadRequest(w, "Error getting scale value")
	}
	json.NewEncoder(w).Encode(&Response{
		Status: "success",
		Data:   val,
	})
})

// executeJob from the requested service
var executeJob = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status string `json:"status"`
	}
	w.Header().Set("Content-Type", "application/json")
	var id int64
	params := mux.Vars(r)
	service, ok := params["service"]
	if !ok {
		log.Error("No service specified")
		handleBadRequest(w, "No service specified")
		return
	}

	idParam, ok := params["id"]
	if !ok {
		log.Error("No ID specified")
		handleBadRequest(w, "No ID specified")
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 0)
	if err != nil {
		log.Error("Error converting (" + idParam + ") to int")
		log.Error(err.Error())
		handleBadRequest(w, "Error converting ("+idParam+") to int")
		return
	}

	q, ok := TriggerQueue[service]
	if !ok {
		log.Error("Invalid Service: " + service)
		handleBadRequest(w, "Invalid Service: "+service)
	}

	elem, err := q.GetByID(id)
	if err != nil {
		log.Error("Could not find job with ID " + idParam)
		handleBadRequest(w, "Could not find job with ID "+idParam)
		return
	}

	// Execute the action
	executeAction(elem)
	json.NewEncoder(w).Encode(&Response{
		Status: "success",
	})
})

// deleteJob from the requested service
var deleteJob = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status string `json:"status"`
	}
	w.Header().Set("Content-Type", "application/json")
	var status string
	var id int64
	params := mux.Vars(r)
	service, ok := params["service"]
	if !ok {
		log.Error("No service specified")
		handleBadRequest(w, "No service specified")
		return
	}

	idParam, ok := params["id"]
	if !ok {
		log.Error("No ID specified")
		handleBadRequest(w, "No ID specified")
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 0)
	if err != nil {
		log.Error("Error converting (" + idParam + ") to int")
		log.Error(err.Error())
		handleBadRequest(w, "Error converting ("+idParam+") to int")
		return
	}

	q, ok := TriggerQueue[service]
	if !ok {
		log.Error("Invalid Service: " + service)
		handleBadRequest(w, "Invalid Service: "+service)
	}

	ok = q.RemoveID(id)
	if !ok {
		log.Error("Could not remove job with ID " + idParam)
		handleBadRequest(w, "Could not remove job with ID "+idParam)
		return
	}
	json.NewEncoder(w).Encode(&Response{
		Status: status,
	})
})

// clearQueue of a particular service
var clearQueue = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type Response struct {
		Status string `json:"status"`
	}
	params := mux.Vars(r)
	service, ok := params["service"]
	if !ok {
		log.Error("No service specified")
		handleBadRequest(w, "No service specified")
		return
	}

	q, ok := TriggerQueue[service]

	if !ok {
		log.Error("Invalid Service: " + service)
		handleBadRequest(w, "Invalid Service: "+service)
		return
	}

	q.Clean()
	json.NewEncoder(w).Encode(&Response{
		Status: "success",
	})
})

// Show Queue(s)
var showQueue = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	service, ok := mux.Vars(r)["service"]
	if !ok {
		log.Error("No service specified")
		handleBadRequest(w, "No service specified")
		return
	}

	q, ok := TriggerQueue[service]
	if !ok {
		log.Error("Invalid Service: " + service)
		handleBadRequest(w, "Invalid Service: "+service)
		return
	}
	json.NewEncoder(w).Encode(q)
})

// Encode all queues
var showAllQueues = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type Response struct {
		Status string                  `json:"status"`
		Body   map[string]*queue.Queue `json:"body"`
	}

	body := make(map[string]*queue.Queue)
	for k, q := range TriggerQueue {
		body[k] = q
	}

	json.NewEncoder(w).Encode(&Response{
		Status: "success",
		Body:   body,
	})
})

// Server testing controllers
var Pong = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	log.Info("Pong!")
	log.Info("%+v", r.Header)
	w.Write([]byte("Pong!\n"))
})
