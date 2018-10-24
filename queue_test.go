package main

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test Queue Functionality
func TestQueueUtils(t *testing.T) {
	q := &Queue{nodes: list.New()}
	head := &Node{1}
	body := &Node{2}
	tail := &Node{3}

	// Enqueue items
	q.Enqueue(head)
	q.Enqueue(body)
	q.Enqueue(tail)

	// Check Queue Util Functions
	assert.Equal(t, q.Head().(*Node), head, "Queue Head is not the correct head")
	assert.Equal(t, q.Tail().(*Node), tail, "Queue Tail is not the correct tail")
	assert.Equal(t, q.Count(), 3, "Queue Count is not the correct count")

	// Dequeue and validate
	assert.Equal(t, q.Dequeue().(*Node), head)
	assert.Equal(t, q.Dequeue().(*Node), body)
	assert.Equal(t, q.Dequeue().(*Node), tail)
}
