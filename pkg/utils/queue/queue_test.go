package queue

import (
	"testing"
)

func TestQueue_Basic(t *testing.T) {
	q := New[int]()
	if !q.IsEmpty() {
		t.Error("queue should be empty initially")
	}
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	if q.Size() != 3 {
		t.Errorf("expected size 3, got %d", q.Size())
	}
	h, err := q.PickHead()
	if err != nil || h != 1 {
		t.Errorf("expected head 1, got %v, err: %v", h, err)
	}
	tail, err := q.PickTail()
	if err != nil || tail != 3 {
		t.Errorf("expected tail 3, got %v, err: %v", tail, err)
	}
	v, err := q.Dequeue()
	if err != nil || v != 1 {
		t.Errorf("expected dequeue 1, got %v, err: %v", v, err)
	}
	v, err = q.Dequeue()
	if err != nil || v != 2 {
		t.Errorf("expected dequeue 2, got %v, err: %v", v, err)
	}
	v, err = q.Dequeue()
	if err != nil || v != 3 {
		t.Errorf("expected dequeue 3, got %v, err: %v", v, err)
	}
	if !q.IsEmpty() {
		t.Error("queue should be empty after all dequeues")
	}
	_, err = q.Dequeue()
	if err == nil {
		t.Error("expected error on dequeue from empty queue")
	}
}

func TestQueue_Iterator(t *testing.T) {
	q := New[int]()
	values := []int{10, 20, 30, 40, 50}
	for _, v := range values {
		q.Enqueue(v)
	}

	var iterated []int
	for v := range q.Iterator() {
		iterated = append(iterated, v)
	}

	if len(iterated) != len(values) {
		t.Fatalf("expected %d elements, got %d", len(values), len(iterated))
	}
	for i, v := range values {
		if iterated[i] != v {
			t.Errorf("at index %d: expected %d, got %d", i, v, iterated[i])
		}
	}

	// After iteration, queue should be empty
	if !q.IsEmpty() {
		t.Errorf("queue should be empty after iteration")
	}
}
