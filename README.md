# Trigger-Queue

## Setup

Ensure you have a go version compatible with modules

```bash
git clone https://github.com/Smart-Chef/trigger-queue # could also probably go get
cd trigger-queue
go build
./trigger-queue
```

## Dev Usage
### Test (ping-pong)
```bash
curl http://localhost:8000/api/ping # response is pong
```

### Add Job to NLP Queue
```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"service":"walk-through","action_key":"sendToNLP","action_params":"test action params","trigger_keys":["temp_>"],"trigger_params": [300]}' \
  http://localhost:8000/api/add/nlp
```

### Delete Job From NLP Queue
```bash
curl --header "Content-Type: application/json" \
  --request POST \
  http://localhost:8000/api/delete/nlp/{id}
```

### Clear NLP Queue
```bash
curl --header "Content-Type: application/json" \
  --request POST \
  http://localhost:8000/api/clear/nlp
```

### Show NLP Queue
```bash
curl --header "Content-Type: application/json" \
  --request POST \
  http://localhost:8000/api/show/nlp
```

### Show All Queues
```bash
curl --header "Content-Type: application/json" \
  --request POST \
  http://localhost:8000/api/show
```
