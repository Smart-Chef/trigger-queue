package queue

import (
	"encoding/json"
	"errors"
	"math/rand"
	"sync"
)

const minQueueLen = 32

type Queue struct {
	items             map[int64]interface{}
	ids               map[interface{}]int64
	buf               []int64
	head, tail, count int
	mutex             *sync.Mutex
	notEmpty          *sync.Cond
	NotEmpty          chan struct{}
}

func (q *Queue) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Items map[int64]interface{} `json:"items"`
	}{
		Items: q.items,
	})
}

func New() *Queue {
	q := &Queue{
		items:    make(map[int64]interface{}),
		ids:      make(map[interface{}]int64),
		buf:      make([]int64, minQueueLen),
		mutex:    &sync.Mutex{},
		NotEmpty: make(chan struct{}, 1),
	}

	q.notEmpty = sync.NewCond(q.mutex)

	return q
}

// Removes all elements from queue
func (q *Queue) Clean() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.items = make(map[int64]interface{})
	q.ids = make(map[interface{}]int64)
	q.buf = make([]int64, minQueueLen)
}

// Returns the number of elements in queue
func (q *Queue) Length() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	return len(q.items)
}

// resize the queue to fit exactly twice its current contents
// this can result in shrinking if the queue is less than half-full
func (q *Queue) resize() {
	newCount := q.count << 1

	if q.count < 2<<18 {
		newCount = newCount << 2
	}

	newBuf := make([]int64, newCount)

	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		n := copy(newBuf, q.buf[q.head:])
		copy(newBuf[n:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

func (q *Queue) notify() {
	if len(q.items) > 0 {
		select {
		case q.NotEmpty <- struct{}{}:
		default:
		}
	}
}

// append util without locking
func (q *Queue) append(elem interface{}, id int64) {
	if q.count == len(q.buf) {
		q.resize()
	}

	q.items[id] = elem
	q.ids[elem] = id
	q.buf[q.tail] = id
	// bitwise modulus
	q.tail = (q.tail + 1) & (len(q.buf) - 1)
	q.count++

	q.notify()

	if q.count == 1 {
		q.notEmpty.Broadcast()
	}
}

// Adds one element at the back of the queue
func (q *Queue) Append(elem interface{}) int64 {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	id := q.newId()
	q.append(elem, id)
	return id
}

func (q *Queue) newId() int64 {
	for {
		id := rand.Int63()
		_, ok := q.items[id]
		if id != 0 && !ok {
			return id
		}
	}
}

// Adds one element at the front of queue
func (q *Queue) Prepend(elem interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if q.count == len(q.buf) {
		q.resize()
	}

	q.head = (q.head - 1) & (len(q.buf) - 1)
	id := q.newId()
	q.items[id] = elem
	q.ids[elem] = id
	q.buf[q.head] = id
	// bitwise modulus
	q.count++

	q.notify()

	if q.count == 1 {
		q.notEmpty.Broadcast()
	}
}

// Previews element at the front of queue
func (q *Queue) Front() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	id := q.buf[q.head]
	if id != 0 {
		return q.items[id]
	}
	return nil
}

// Evaluate Front Element element at the front of queue for trigger queue
func (q *Queue) EvaluateFront(triggerFunc func(interface{}) (bool, []error), onTrigErr func(interface{}, []error)) (bool, *interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Return if the queue is empty
	if q.count <= 0 {
		return false, nil
	}

	// Get the id of our head and pop it from the buffer
	id := q.pop()

	// Ensure that our buffer is not empty -- maybe?
	if id != 0 {
		// Get the head element
		elem := q.items[id]

		// Remove the element from the head of the ids and items list
		delete(q.ids, elem)
		delete(q.items, id)

		// Check if all the triggers evaluate to true
		success, err := triggerFunc(elem)
		if err != nil {
			onTrigErr(elem, err)
		}
		if success {
			return success, &elem
		}

		// Append the item to the end of the queue
		q.append(elem, id)
	}
	return false, nil
}

// Previews element at the back of queue
func (q *Queue) Back() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	id := q.buf[(q.tail-1)&(len(q.buf)-1)]
	if id != 0 {
		return q.items[id]
	}
	return nil
}

func (q *Queue) pop() int64 {
	for {
		if q.count <= 0 {
			q.notEmpty.Wait()
		}

		// I have no idea why, but sometimes it's less than 0
		if q.count > 0 {
			break
		}
	}

	id := q.buf[q.head]
	q.buf[q.head] = 0

	// bitwise modulus
	q.head = (q.head + 1) & (len(q.buf) - 1)
	q.count--
	if len(q.buf) > minQueueLen && (q.count<<1) == len(q.buf) {
		q.resize()
	}

	return id
}

// Pop removes and returns the element from the front of the queue.
// If the queue is empty, it will block
func (q *Queue) Pop() interface{} {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for {
		id := q.pop()

		item, ok := q.items[id]

		if ok {
			delete(q.ids, item)
			delete(q.items, id)
			q.notify()
			return item
		}
	}
}

// Removes one element from the queue
func (q *Queue) Remove(elem interface{}) bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	id, ok := q.ids[elem]
	if !ok {
		return false
	}
	delete(q.ids, elem)
	delete(q.items, id)
	return true
}

// Removes one element from the queue by ID
func (q *Queue) RemoveID(id int64) bool {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	elem, ok := q.items[id]

	if !ok {
		return false
	}

	delete(q.ids, elem)
	delete(q.items, id)
	return true
}

// Get a single element by ID also removing it from the queue
func (q *Queue) GetByID(id int64) (interface{}, error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	elem, ok := q.items[id]

	if !ok {
		return nil, errors.New("no element with specified id")
	}

	delete(q.ids, elem)
	delete(q.items, id)
	return elem, nil
}
