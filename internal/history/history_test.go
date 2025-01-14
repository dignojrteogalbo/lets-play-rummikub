package history

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStack(t *testing.T) {
	newStack := NewStack[int]()
	assert.NotNil(t, newStack)
	assert.IsType(t, ([]int)(nil), newStack.(*stack[int]).states)
}

func TestPop(t *testing.T) {
	states := []string{"hello", "world", "1", "2"}
	newStack := &stack[string]{states}
	assert.Len(t, newStack.states, 4)
	popped := newStack.Pop()
	assert.Len(t, newStack.states, 3)
	assert.Equal(t, popped, "2")
	popped = newStack.Pop()
	assert.Len(t, newStack.states, 2)
	assert.Equal(t, popped, "1")
	popped = newStack.Pop()
	assert.Len(t, newStack.states, 1)
	assert.Equal(t, popped, "world")
	popped = newStack.Pop()
	assert.Equal(t, popped, "hello")
	assert.Empty(t, newStack.states)
	popped = newStack.Pop()
	assert.Equal(t, popped, "")
}

func TestPush(t *testing.T) {
	newStack := &stack[string]{states: make([]string, 0)}
	newStack.Push("bottomOfStack")
	assert.Len(t, newStack.states, 1)
	newStack.Push("afterBottom")
	assert.Len(t, newStack.states, 2)
	newStack.Push("beforeTop")
	assert.Len(t, newStack.states, 3)
	newStack.Push("topOfStack")
	assert.Len(t, newStack.states, 4)
	assert.Equal(t, newStack.states[0], "bottomOfStack")
	assert.Equal(t, newStack.states[1], "afterBottom")
	assert.Equal(t, newStack.states[2], "beforeTop")
	assert.Equal(t, newStack.states[3], "topOfStack")
}

func TestClear(t *testing.T) {
	states := []string{"hello", "world", "1", "2"}
	history := &stack[string]{states}
	history.Clear()
	assert.Empty(t, history.states)
}
