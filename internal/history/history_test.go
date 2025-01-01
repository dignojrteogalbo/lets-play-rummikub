package history

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCloneable struct {
	value string
}

func (m *mockCloneable) Clone() Cloneable {
	return &mockCloneable{m.value}
}

func (m *mockCloneable) Restore(restore Cloneable) {
	rep, ok := restore.(*mockCloneable)
	if !ok {
		return
	}
	m.value = rep.value
}

func TestNewHistory(t *testing.T) {
	original := &mockCloneable{"hello world"}
	history := NewHistory(original)
	assert.NotNil(t, history)
	assert.IsType(t, (*historyInstance)(nil), history)
	assert.Same(t, history.(*historyInstance).original, original)
}

func TestUndo(t *testing.T) {
	states := []Cloneable{
		&mockCloneable{"hello"},
		&mockCloneable{"world"},
		&mockCloneable{"1"},
		&mockCloneable{"2"},
	}
	current := &mockCloneable{"3"}
	history := &historyInstance{original: current, states: make([]Cloneable, len(states))}
	copy(history.states, states)
	recent := history.Undo()
	assert.Same(t, recent, states[3])
	assert.NotSame(t, current, recent)
	assert.Equal(t, current, recent)
	recent = history.Undo()
	assert.Same(t, recent, states[2])
	assert.Equal(t, current.value, "1")
	recent = history.Undo()
	assert.Same(t, recent, states[1])
	assert.Equal(t, current.value, "world")
	recent = history.Undo()
	assert.Same(t, recent, states[0])
	assert.Equal(t, current.value, "hello")
	empty := history.Undo()
	assert.Nil(t, empty)
	assert.Equal(t, current.value, "hello")
}

func TestAppend(t *testing.T) {
	original := &mockCloneable{"start"}
	history := NewHistory(original).(*historyInstance)
	bottomOfStack := &mockCloneable{"1"}
	afterBottom := &mockCloneable{"1"}
	beforeTop := &mockCloneable{"1"}
	topOfStack := &mockCloneable{"1"}
	history.Append(bottomOfStack)
	assert.Len(t, history.states, 1)
	history.Append(afterBottom)
	assert.Len(t, history.states, 2)
	history.Append(beforeTop)
	assert.Len(t, history.states, 3)
	history.Append(topOfStack)
	assert.Len(t, history.states, 4)
	assert.Same(t, history.states[0], bottomOfStack)
	assert.Same(t, history.states[1], afterBottom)
	assert.Same(t, history.states[2], beforeTop)
	assert.Same(t, history.states[3], topOfStack)
}

func TestClear(t *testing.T) {
	states := []Cloneable{
		&mockCloneable{"hello"},
		&mockCloneable{"world"},
		&mockCloneable{"1"},
		&mockCloneable{"2"},
	}
	current := &mockCloneable{"3"}
	history := &historyInstance{original: current, states: make([]Cloneable, len(states))}
	history.Clear()
	assert.Empty(t, history.states)
}
