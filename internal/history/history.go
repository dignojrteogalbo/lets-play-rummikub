package history

type (
	Undoable interface {
		Undo()
	}

	Stack[T any] interface {
		Pop() T
		Push(T)
		Clear()
	}

	stack[T any] struct {
		states []T
	}
)

func (s *stack[T]) Pop() T {
	if len(s.states) == 0 {
		return *new(T)
	}
	top := s.states[len(s.states)-1]
	s.states = s.states[:len(s.states)-1]
	return top
}

func (h *stack[T]) Push(push T) {
	h.states = append(h.states, push)
}

func (h *stack[T]) Clear() {
	h.states = nil
}

func NewStack[T any]() Stack[T] {
	return &stack[T]{states: make([]T, 0)}
}
