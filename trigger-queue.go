package main

import (
	"encoding/json"
	"errors"
	"time"
	"trigger-queue/queue"

	log "github.com/sirupsen/logrus"
)

var TriggerQueue = map[string]*queue.Queue{
	"nlp":          queue.New(),
	"walk-through": queue.New(),
	"other":        queue.New(),
}

type Job struct {
	ID            int64 `json:"id"`
	Triggers      []Trigger
	TriggerKeys   []string `json:"trigger_keys"`
	TriggerParams []interface{}
	Action        Action
	ActionKey     string
	ActionParams  interface{}
	Subscriber    string
	CreatedAt     time.Time
	Errors        []error
}

func (j *Job) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID            int64         `json:"id"`
		TriggerKeys   []string      `json:"trigger_keys"`
		TriggerParams []interface{} `json:"trigger_params"`
		ActionKey     string        `json:"action_key"`
		ActionParams  interface{}   `json:"action_params"`
		Subscriber    string        `json:"subscriber"`
		CreatedAt     time.Time     `json:"created_at"`
		Errors        []string      `json:"errors"`
	}{
		ID:            j.ID,
		TriggerKeys:   j.TriggerKeys,
		TriggerParams: j.TriggerParams,
		ActionKey:     j.ActionKey,
		ActionParams:  j.ActionParams,
		Subscriber:    j.Subscriber,
		CreatedAt:     j.CreatedAt,
		Errors:        serializeListErrs(j.Errors),
	})
}

func serializeListErrs(errs []error) []string {
	s := make([]string, 0)
	for _, e := range errs {
		s = append(s, e.Error())
	}
	return s
}

func createJob(subscriber string, triggerKeys []string, triggerParams []interface{}, actionKey string, actionParams interface{}) (*Job, error) {
	var ok bool
	var action Action
	var trigger Trigger
	triggers := make([]Trigger, len(triggerKeys))

	// Get Action
	action, ok = Actions[actionKey]

	if !ok {
		return nil, errors.New("No action found named \"" + actionKey + "\"")
	}

	// Get triggers -- should no triggerKeys be acceptable?
	for i, key := range triggerKeys {
		trigger, ok = Triggers[key]
		if !ok {
			return nil, errors.New("No trigger found named \"" + key + "\"")
		}
		triggers[i] = trigger
	}

	j := &Job{
		Triggers:      triggers,
		TriggerKeys:   triggerKeys,
		TriggerParams: triggerParams,
		Action:        action,
		ActionKey:     actionKey,
		ActionParams:  actionParams,
		Subscriber:    subscriber,
		CreatedAt:     time.Now(),
	}

	out, _ := json.Marshal(j)
	log.Info("Created Job", string(out))
	return j, nil
}

func evaluateTriggers(j interface{}) (bool, []error) {
	if j == nil {
		return false, nil
	}
	var trigger Trigger
	var triggerParam interface{}
	errs := make([]error, 0)
	job := j.(*Job)
	i := 0
	passed := true

	// ignore if the job has any existing errors
	if len(job.Errors) > 0 {
		return false, job.Errors
	}

	// Evaluate all the triggers
	for passed && i < len(job.Triggers) {
		trigger = job.Triggers[i]
		triggerParam = job.TriggerParams[i]
		p, err := trigger(triggerParam)
		if err != nil {
			errs = append(errs, err)
		}
		passed = p
		i++
	}
	return passed, errs
}

func executeAction(j interface{}) {
	job := j.(*Job)
	log.Printf("Executing Action (%s)", job.ActionKey)
	job.Action(job.ActionParams)
}

func triggerError(j interface{}, errs []error) {
	j.(*Job).Errors = errs
}
