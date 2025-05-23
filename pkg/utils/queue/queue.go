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

func New[T any]() Queue[T] {
	return Queue[T]{
		head: nil,
		tail: nil,
		size: 0,
	}
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
