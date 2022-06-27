package common

type (
	stackNode struct {
		value any
		next  *stackNode
	}
	Stack[T any] struct {
		top *stackNode
	}
)

func (s *Stack[T]) Push(value T) {
	s.top = &stackNode{
		value: value,
		next:  s.top,
	}
}

func (s *Stack[T]) Peek() T {
	return s.top.value.(T)
}

func (s *Stack[T]) Pop() T {
	value := s.top.value.(T)
	s.top = s.top.next
	return value
}
