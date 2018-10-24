package main

import "time"

type Job struct {
	triggers      []Trigger
	triggerParams interface{}
	action        Action
	actionParams  interface{}
	subscriber    string
	createdAt     time.Time
}
