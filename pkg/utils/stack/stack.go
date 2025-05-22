package stack

import "fmt"

type (
	Stack[T any] struct {
		top *node[T]
	}
	node[T any] struct {
		value T
		prev  *node[T]
	}
)

func New[T any]() Stack[T] {
	return Stack[T]{
		top: nil,
	}
}

func (s *Stack[T]) Push(v T) {
	node := &node[T]{
		value: v,
		prev:  s.top,
	}
	s.top = node
}

func (s *Stack[T]) Pop() (T, error) {
	if s.top == nil {
		var zero T
		return zero, fmt.Errorf("stack is empty")
	}
	node := s.top
	s.top = node.prev

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
	return s.top == nil
}
