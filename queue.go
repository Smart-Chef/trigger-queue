package main

import "container/list"

type Node struct {
	Value int
}

// Queue is a basic FIFO queue based on a doubly linked list
type Queue struct {
	nodes *list.List
}

// Push adds a node to the queue.
func (q *Queue) Enqueue(n interface{}) {
	q.nodes.PushBack(n)
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) Dequeue() interface{} {
	return q.nodes.Remove(q.nodes.Front())
}

func (q *Queue) Count() int {
	return q.nodes.Len()
}

func (q *Queue) Head() interface{} {
	return q.nodes.Front().Value
}

func (q *Queue) Tail() interface{} {
	return q.nodes.Back().Value
}
