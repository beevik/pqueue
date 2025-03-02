package pq

import (
	"cmp"
)

// item is a struct holding the key and value for an item within the priority
// queue.
type item[K cmp.Ordered, V any] struct {
	Key   K
	Value V
}

// Queue implements a min-heap priority queue.
type Queue[K cmp.Ordered, V any] struct {
	items []item[K, V]
}

// NewQueue creates a new priority queue.
func NewQueue[K cmp.Ordered, V any]() *Queue[K, V] {
	return &Queue[K, V]{
		items: make([]item[K, V], 0),
	}
}

// Count returns the number of elements in the priority queue.
func (q *Queue[K, V]) Count() int {
	return len(q.items)
}

// IsEmpty returns true if the priority queue is empty.
func (q *Queue[K, V]) IsEmpty() bool {
	return len(q.items) == 0
}

// Enqueue adds an item to the priority queue.
func (q *Queue[K, V]) Enqueue(key K, value V) {
	item := item[K, V]{key, value}
	q.items = append(q.items, item)
	q.siftUp(len(q.items) - 1)
}

// Dequeue removes and returns the item with the lowest priority. Panics
// if called on an empty queue.
func (q *Queue[K, V]) Dequeue() (K, V) {
	if len(q.items) == 0 {
		panic("Dequeue performed on an empty priority queue")
	}

	min := q.items[0]
	last := len(q.items) - 1
	q.items[0] = q.items[last]
	q.items = q.items[:last]

	if len(q.items) > 0 {
		q.siftDown(0)
	}

	return min.Key, min.Value
}

// TryDequeue removes and returns the item with the lowest priority. If the
// queue was empty, TryDequeue returns the zero key, zero value and false.
func (q *Queue[K, V]) TryDequeue() (K, V, bool) {
	if len(q.items) == 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	k, v := q.Dequeue()
	return k, v, true
}

// Peek returns the item with the lowest priority without removing it.
// Panics if called on an empty queue.
func (q *Queue[K, V]) Peek() (K, V) {
	if len(q.items) == 0 {
		panic("Peek performed on an empty priority queue")
	}

	min := q.items[0]
	return min.Key, min.Value
}

// TryPeek returns the item with the lowest priority without removing it. If
// the queue is empty, TryPeek returns a zero key, zero value, and false.
func (q *Queue[K, V]) TryPeek() (K, V, bool) {
	if len(q.items) == 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	k, v := q.Peek()
	return k, v, true
}

// siftUp moves the element at index i up to its correct position.
func (q *Queue[K, V]) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if q.items[parent].Key <= q.items[i].Key {
			break
		}
		q.items[parent], q.items[i] = q.items[i], q.items[parent]
		i = parent
	}
}

// siftDown moves the element at index i down to its correct position.
func (q *Queue[K, V]) siftDown(i int) {
	n := len(q.items)
	for {
		min := i
		left := 2*i + 1
		right := 2*i + 2

		if left < n && q.items[left].Key < q.items[min].Key {
			min = left
		}

		if right < n && q.items[right].Key < q.items[min].Key {
			min = right
		}

		if min == i {
			break
		}

		q.items[i], q.items[min] = q.items[min], q.items[i]
		i = min
	}
}
