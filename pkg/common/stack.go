package common

type (
	stackNode struct {
		Value any
		Next  *stackNode
	}
	Stack[T any] struct {
		Top *stackNode
	}
)

func (s *Stack[T]) Push(value T) {
	s.Top = &stackNode{
		Value: value,
		Next:  s.Top,
	}
}

func (s *Stack[T]) Peek() T {
	return s.Top.Value.(T)
}

func (s *Stack[T]) Pop() T {
	value := s.Top.Value.(T)
	s.Top = s.Top.Next
	return value
}

func (s *Stack[T]) HasNext() bool {
	return s.Top != nil
}
