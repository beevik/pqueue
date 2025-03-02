package pq

import (
	"testing"
)

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	f()
}

func TestNewQueue(t *testing.T) {
	q := NewQueue[int, string]()
	if q == nil {
		t.Fatal("NewQueue returned nil")
	}
	if q.Count() != 0 {
		t.Errorf("New queue should be empty, got count: %d", q.Count())
	}
	if !q.IsEmpty() {
		t.Errorf("New queue should have been empty")
	}
	if len(q.items) != 0 {
		t.Errorf("New queue should have empty items slice, got length: %d", len(q.items))
	}
}

func TestEnqueue(t *testing.T) {
	q := NewQueue[int, string]()

	// Test enqueueing a single item
	q.Enqueue(5, "five")
	if q.Count() != 1 {
		t.Errorf("Queue should have 1 item, got: %d", q.Count())
	}
	if q.items[0].Key != 5 || q.items[0].Value != "five" {
		t.Errorf("Expected (5, five), got (%v, %v)", q.items[0].Key, q.items[0].Value)
	}

	// Test enqueueing multiple items
	q.Enqueue(3, "three")
	q.Enqueue(7, "seven")

	if q.Count() != 3 {
		t.Errorf("Queue should have 3 items, got: %d", q.Count())
	}

	// Check heap property
	if q.items[0].Key != 3 {
		t.Errorf("Min element should be 3, got %v", q.items[0].Key)
	}
}

func TestPeek(t *testing.T) {
	q := NewQueue[int, string]()

	// Test peek on empty queue
	assertPanic(t, func() {
		q.Peek()
	})

	// Test peek with items
	q.Enqueue(5, "five")
	q.Enqueue(3, "three")
	q.Enqueue(7, "seven")

	k, v := q.Peek()
	if k != 3 || v != "three" {
		t.Errorf("Expected peek to return (3, three), got (%v, %v)", k, v)
	}

	var ok bool
	k, v, ok = q.TryPeek()
	if k != 3 || v != "three" || !ok {
		t.Errorf("Expected TryPeek tro return (3, three, true), got (%v, %v, %v)", k, v, ok)
	}

	// Ensure queue size didn't change
	if q.Count() != 3 {
		t.Errorf("Queue size should remain 3 after peek, got: %d", q.Count())
	}
}

func TestDequeue(t *testing.T) {
	q := NewQueue[int, string]()

	// Test empty dequeue
	assertPanic(t, func() {
		q.Dequeue()
	})

	// Add items in non-sorted order
	q.Enqueue(5, "five")
	q.Enqueue(3, "three")
	q.Enqueue(7, "seven")
	q.Enqueue(1, "one")
	q.Enqueue(9, "nine")

	// Verify dequeue gets min element
	k, v, ok := q.TryDequeue()
	if k != 1 || v != "one" || !ok {
		t.Errorf("Expected dequeue to return (1, one, true), got (%v, %v, %v)", k, v, ok)
	}
	if q.Count() != 4 {
		t.Errorf("Queue size should be 4 after dequeue, got: %d", q.Count())
	}

	// Verify next dequeue
	k, v = q.Dequeue()
	if k != 3 || v != "three" {
		t.Errorf("Expected dequeue to return (3, three), got (%v, %v)", k, v)
	}

	// Dequeue remaining items
	q.Dequeue()        // 5
	q.Dequeue()        // 7
	k, v = q.Dequeue() // 9
	if k != 9 || v != "nine" {
		t.Errorf("Expected dequeue to return (9, nine), got (%v, %v)", k, v)
	}

	// Queue should now be empty
	if q.Count() != 0 {
		t.Errorf("Queue should be empty after dequeueing all items, got count: %d", q.Count())
	}
}

func TestWithDifferentTypes(t *testing.T) {
	// Test with float keys
	qFloat := NewQueue[float64, int]()
	qFloat.Enqueue(3.14, 314)
	qFloat.Enqueue(1.618, 1618)
	qFloat.Enqueue(2.718, 2718)

	k, v := qFloat.Dequeue()
	if k != 1.618 || v != 1618 {
		t.Errorf("Expected (1.618, 1618), got (%v, %v)", k, v)
	}

	// Test with string keys
	qString := NewQueue[string, bool]()
	qString.Enqueue("banana", true)
	qString.Enqueue("apple", false)
	qString.Enqueue("cherry", true)

	k2, v2 := qString.Dequeue()
	if k2 != "apple" || v2 != false {
		t.Errorf("Expected (apple, false), got (%v, %v)", k2, v2)
	}
}

func TestWithEqualPriorities(t *testing.T) {
	q := NewQueue[int, string]()

	// Add items with equal priorities
	q.Enqueue(5, "first-five")
	q.Enqueue(5, "second-five")
	q.Enqueue(5, "third-five")

	if q.Count() != 3 {
		t.Errorf("Queue should have 3 items, got: %d", q.Count())
	}

	// Dequeue should return one of the items with priority 5
	k, _ := q.Dequeue()
	if k != 5 {
		t.Errorf("Expected key 5, got %v", k)
	}

	// Should have 2 items left
	if q.Count() != 2 {
		t.Errorf("Queue should have 2 items left, got: %d", q.Count())
	}

	// Check that all returned items have key 5
	k, _ = q.Dequeue()
	if k != 5 {
		t.Errorf("Expected key 5, got %v", k)
	}

	k, _ = q.Dequeue()
	if k != 5 {
		t.Errorf("Expected key 5, got %v", k)
	}

	// Queue should now be empty
	if q.Count() != 0 {
		t.Errorf("Queue should be empty, got count: %d", q.Count())
	}
}

func TestHeapOperations(t *testing.T) {
	// This test verifies that the heap operations work correctly
	// by enqueueing items in non-sorted order and checking that
	// they're dequeued in the correct order
	q := NewQueue[int, string]()

	// Add items in random order
	testData := []struct {
		key   int
		value string
	}{
		{10, "ten"},
		{5, "five"},
		{8, "eight"},
		{3, "three"},
		{7, "seven"},
		{1, "one"},
		{4, "four"},
		{9, "nine"},
		{2, "two"},
		{6, "six"},
	}

	for _, item := range testData {
		q.Enqueue(item.key, item.value)
	}

	// Expected dequeue order
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Dequeue all items and verify they come out in sorted order
	for i := 0; i < len(expected); i++ {
		k, _ := q.Dequeue()
		if k != expected[i] {
			t.Errorf("Expected %d at position %d, got %d", expected[i], i, k)
		}
	}
}

func TestEmptyQueue(t *testing.T) {
	q := NewQueue[int, string]()

	// Verify Count returns 0
	if q.Count() != 0 {
		t.Errorf("Empty queue should have count 0, got: %d", q.Count())
	}

	// Verify Peek panics on an empty queue
	assertPanic(t, func() {
		q.Peek()
	})

	// Verify Dequeue panics on an empty queue
	assertPanic(t, func() {
		q.Dequeue()
	})

	// Verify TryDequeue returns ok==false.
	_, _, ok := q.TryDequeue()
	if ok {
		t.Errorf("Queue should have been empty")
	}

	// Verify TryPeek returns ok==false
	_, _, ok = q.TryPeek()
	if ok {
		t.Errorf("Queue should have been empty")
	}

	// Enqueue one item then dequeue it to make queue empty again
	q.Enqueue(1, "one")
	q.Dequeue()

	// Verify Count returns 0 again
	if q.Count() != 0 {
		t.Errorf("Queue should be empty after dequeueing all items, got count: %d", q.Count())
	}
}
