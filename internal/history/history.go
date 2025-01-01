package history

type (
	Cloneable interface {
		Clone() Cloneable
		Restore(Cloneable)
	}

	History interface {
		Undo() Cloneable
		Append(Cloneable)
		Clear()
	}

	historyInstance struct {
		original Cloneable
		states   []Cloneable
	}
)

func (h *historyInstance) Undo() Cloneable {
	if len(h.states) == 0 {
		return nil
	}
	top := h.states[len(h.states)-1]
	h.original.Restore(top)
	h.states = h.states[:len(h.states)-1]
	return top
}

func (h *historyInstance) Append(clone Cloneable) {
	h.states = append(h.states, clone)
}

func (h *historyInstance) Clear() {
	h.states = nil
}

func NewHistory(original Cloneable) History {
	return &historyInstance{original: original}
}
