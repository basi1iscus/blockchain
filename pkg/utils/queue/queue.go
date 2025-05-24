package queue

import "fmt"

type (
	Queue[T any] struct {
		head *node[T]
		tail *node[T]
		size int
	}
	node[T any] struct {
		value T
		prev  *node[T]
		next  *node[T]
	}
)

func New[T any]() *Queue[T] {
	return &Queue[T]{
		head: nil,
		tail: nil,
		size: 0,
	}
}

func FromArray[T any](items []T) *Queue[T] {
	q := New[T]()
	for _, item := range items {
		q.Enqueue(item)
	}

	return q
}

func (s *Queue[T]) ToArray() []T {
	arr := make([]T, 0, s.size)
	for node := s.head; node != nil; node = node.next {
		arr = append(arr, node.value)
	}
	return arr
}

func (s *Queue[T]) Enqueue(v T) {
	node := &node[T]{
		value: v,
		prev:  s.tail,
	}
	if s.tail == nil {
		s.head = node
	} else {
		s.tail.next = node
	}
	s.tail = node
	s.size++
}

func (s *Queue[T]) Dequeue() (T, error) {
	if s.head == nil {
		var zero T
		return zero, fmt.Errorf("queue is empty")
	}
	node := s.head
	s.head = node.next
	if s.head == nil {
		s.tail = nil
	} else {
		s.head.prev = nil
	}
	s.size--

	return node.value, nil
}

func (s *Queue[T]) PickTail() (T, error) {
	if s.tail == nil {
		var zero T
		return zero, fmt.Errorf("queue is empty")
	}
	return s.tail.value, nil
}

func (s *Queue[T]) PickHead() (T, error) {
	if s.head == nil {
		var zero T
		return zero, fmt.Errorf("queue is empty")
	}
	return s.head.value, nil
}

func (s *Queue[T]) IsEmpty() bool {
	return s.size == 0
}

func (s *Queue[T]) Size() int {
	return s.size
}

func (s *Queue[T]) Next() (T, bool) {
    val, err := s.Dequeue()
	if err != nil {	
		 var zero T
		return zero, false
	}
    return val, true
}

func (q *Queue[T]) Iterator() <-chan T {
    ch := make(chan T)
    go func() {
		for v, ok := q.Next(); ok; v, ok = q.Next() {
			ch <-v
		}		

        close(ch)
    }()
    return ch
}