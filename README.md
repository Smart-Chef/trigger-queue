# Trigger-Queue
[![Build Status](https://travis-ci.org/Smart-Chef/trigger-queue.svg?branch=master)](https://travis-ci.org/Smart-Chef/trigger-queue)
[![codecov](https://codecov.io/gh/Smart-Chef/trigger-queue/branch/master/graph/badge.svg)](https://codecov.io/gh/Smart-Chef/trigger-queue)

The trigger queue has two actions and one constantly running event loop. The event loop is constantly dequeue events and evaluating all the triggers. If the triggers do not all return true, the event will be enqueued. If the triggers all return true, the action will then be executed.

The event subscription simply validates the request to subscribe a new request and then enqueus the new event. To remove triggers, a service can pass in some qualifier or filter (usually a uuid). The specified event and all related events will then be removed/unsubscribed from the trigger queue. 

## Setup

Ensure you have a go version compatible with modules

```bash
git clone https://github.com/Smart-Chef/trigger-queue # could also probably go get
cd trigger-queue
go build
./trigger-queue
```

## Dev Usage

The API docs can be found at [PostmanDocs](https://documenter.getpostman.com/view/1907478/RzZ6GzUp)
