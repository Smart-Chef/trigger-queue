package main

import (
	"encoding/json"
	"errors"
	"time"
	"trigger-queue/queue"

	"github.com/sirupsen/logrus"
)

var TriggerQueue = map[string]*queue.Queue{
	"nlp":          queue.New(),
	"walk-through": queue.New(),
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
	}{
		ID:            j.ID,
		TriggerKeys:   j.TriggerKeys,
		TriggerParams: j.TriggerParams,
		ActionKey:     j.ActionKey,
		ActionParams:  j.ActionParams,
		Subscriber:    j.Subscriber,
		CreatedAt:     j.CreatedAt,
	})
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
	logrus.Info("Created Job", out)
	return j, nil
}

func executeTrigger(j interface{}) bool {
	if j == nil {
		return false
	}
	var trigger Trigger
	var triggerParam interface{}
	job := j.(*Job)
	i := 0
	passed := true

	for passed && i < len(job.Triggers) {
		trigger = job.Triggers[i]
		triggerParam = job.TriggerParams[i]
		passed = trigger(triggerParam)
		i++
	}
	return passed
}

func executeAction(j interface{}) {
	logrus.Info("Executing Action")
	job := j.(*Job)
	job.Action(job.ActionParams)
}
