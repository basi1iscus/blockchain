package stack

import (
	"testing"
)

func TestStack_PushPop(t *testing.T) {
	s := New[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)

	v, err := s.Pop()
	if err != nil || v != 3 {
		t.Errorf("Pop() = %v, %v; want 3, nil", v, err)
	}
	v, err = s.Pop()
	if err != nil || v != 2 {
		t.Errorf("Pop() = %v, %v; want 2, nil", v, err)
	}
	v, err = s.Pop()
	if err != nil || v != 1 {
		t.Errorf("Pop() = %v, %v; want 1, nil", v, err)
	}
	_, err = s.Pop()
	if err == nil {
		t.Error("Pop() on empty stack should return error")
	}
}

func TestStack_Pick(t *testing.T) {
	s := New[string]()
	s.Push("a")
	s.Push("b")
	v, err := s.Pick()
	if err != nil || v != "b" {
		t.Errorf("Pick() = %v, %v; want b, nil", v, err)
	}
	s.Pop()
	v, err = s.Pick()
	if err != nil || v != "a" {
		t.Errorf("Pick() = %v, %v; want a, nil", v, err)
	}
	s.Pop()
	_, err = s.Pick()
	if err == nil {
		t.Error("Pick() on empty stack should return error")
	}
}

func TestStack_IsEmpty(t *testing.T) {
	s := New[float64]()
	if !s.IsEmpty() {
		t.Error("IsEmpty() = false; want true")
	}
	s.Push(1.23)
	if s.IsEmpty() {
		t.Error("IsEmpty() = true after Push; want false")
	}
	s.Pop()
	if !s.IsEmpty() {
		t.Error("IsEmpty() = false after Pop; want true")
	}
}
