package stack

import "fmt"

type (
	Stack[T any] struct {
		top *node[T]
		size int
	}
	node[T any] struct {
		value T
		prev  *node[T]
	}
)

func New[T any]() *Stack[T] {
	return &Stack[T]{
		top: nil,
	}
}

func FromArray[T any](items []T) *Stack[T] {
	s := New[T]()
	for _, item := range items {
		s.Push(item)
	}

	return s
}

func (s *Stack[T]) ToArray() []T {
	arr := make([]T, 0, s.size)
	for node := s.top; node != nil; node = node.prev {
		arr = append(arr, node.value)
	}
	return arr
}

func (s *Stack[T]) Push(v T) {
	node := &node[T]{
		value: v,
		prev:  s.top,
	}
	s.top = node
	s.size++
}

func (s *Stack[T]) Pop() (T, error) {
	if s.top == nil {
		var zero T
		return zero, fmt.Errorf("stack is empty")
	}
	node := s.top
	s.top = node.prev
	s.size--

	return node.value, nil
}

func (s *Stack[T]) Pick() (T, error) {
	if s.top == nil {
		var zero T
		return zero, fmt.Errorf("stack is empty")
	}
	return s.top.value, nil
}

func (s *Stack[T]) IsEmpty() bool {
	return s.size == 0
}

func (s *Stack[T]) Size() int {
	return s.size
}
