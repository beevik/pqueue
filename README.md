pq
==

A generic implementation of a priority queue in Go. The queue is implemented
as a min-heap, with the key value used as a priority value. Items with
smallest key value are dequeued first.

```go
q := pq.NewQueue[int, string]()

q.Enqueue(3, "three")
q.Enqueue(1, "one")
q.Enqueue(5, "five")

fmt.Printf("Queue size: %d\n", q.Count())

k, v := q.Peek()
fmt.Printf("Peek: %d %s\n", k, v)

for !q.IsEmpty() {
    k, v = q.Dequeue()
    fmt.Printf("Dequeue: %d %s\n", k, v)
}
```

Output:
```
Queue size: 3
Peek: 1 one
Dequeue: 1 one
Dequeue: 3 three
Dequeue: 5 five
```
