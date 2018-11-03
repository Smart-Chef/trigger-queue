package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"trigger-queue/queue"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

var INVALIDSTATUS = "INVALID"

// handleBadRequest for the application
func handleBadRequest(w http.ResponseWriter, msgAndArgs ...interface{}) {
	//log.Warn(msgAndArgs)
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(&struct {
		Status string      `json:"status"`
		Msg    interface{} `json:"msg"`
	}{INVALIDSTATUS, msgAndArgs})
}

// addJob to the requested service
var addJob = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	params := mux.Vars(r)

	req := &Request{}
	json.NewDecoder(r.Body).Decode(&req)
	w.Header().Set("Content-Type", "application/json")

	q, ok := TriggerQueue[req.Service]
	if !ok {
		handleBadRequest(w, "Unrecognized Service \""+req.Service+"\"")
		return
	}

	j, err := createJob(params["service"], req.TriggerKeys, req.TriggerParams, req.ActionKey, req.ActionParams)
	if err != nil {
		handleBadRequest(w, err.Error())
		return
	}

	id := q.Append(j)
	j.ID = id
	json.NewEncoder(w).Encode(&Response{
		Status: req.Service,
		JobID:  id,
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
	service := params["service"]

	id, err := strconv.ParseInt(params["id"], 10, 0)

	if err != nil {
		handleBadRequest(w, err)
	}

	q := TriggerQueue[service]
	success := q.RemoveID(id)

	if success {
		status = "success"
	} else {
		status = "failure"
	}
	json.NewEncoder(w).Encode(&Response{
		Status: status,
	})
})

// clearQueue of a particular service
var clearQueue = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status string `json:"status"`
	}
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	service := params["service"]

	q, ok := TriggerQueue[service]

	if !ok {
		handleBadRequest(w, "No service \""+service+"\"")
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
	service := mux.Vars(r)["service"]
	q, ok := TriggerQueue[service]

	if !ok {
		handleBadRequest(w, "Invalid service \""+service+"\"")
		return
	}
	json.NewEncoder(w).Encode(q)
})

var showAllQueues = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status string                  `json:"status"`
		Body   map[string]*queue.Queue `json:"body"`
	}
	w.Header().Set("Content-Type", "application/json")
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
	w.Write([]byte("Pong!\n"))
})
