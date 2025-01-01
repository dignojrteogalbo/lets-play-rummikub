package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createRunTiles(t *testing.T, start, end int, color Color) []Piece {
	length := end - start + 1
	assert.True(t, length >= 3)
	tiles := make([]Piece, 0)
	for i := start; i <= end; i++ {
		tiles = append(tiles, NewPiece(Value(i), color))
	}
	return tiles
}

func createGroupTiles(t *testing.T, length uint, value Value) []Piece {
	assert.True(t, length >= 3)
	tiles := make([]Piece, 0)
	for i := uint(0); i < length; i++ {
		color := Color(i%4 + 1)
		tiles = append(tiles, NewPiece(value, color))
	}
	return tiles
}

func TestLen(t *testing.T) {
	nilSet := (*set)(nil)
	assert.Zero(t, nilSet.Len())
	notEmpty := &set{tiles: createRunTiles(t, 1, 9, ColorBlue)}
	assert.Equal(t, notEmpty.Len(), 9)
}

func TestString(t *testing.T) {
	test := &set{tiles: createRunTiles(t, 1, 13, ColorBlack)}
	test.tiles = append(test.tiles, NewPiece(ValueJoker, ColorBlack))
	result := test.String()
	assert.NotEmpty(t, result)
}

func TestClone(t *testing.T) {
	original := &set{tiles: createRunTiles(t, 1, 13, ColorRed)}
	clone := original.Clone()
	assert.Equal(t, original, clone)
	assert.NotSame(t, original, clone.(*set))
}

func TestIsValidSet(t *testing.T) {
	t.Run("ShouldReturnTrueOnGroup", func(t *testing.T) {
		group := &set{tiles: createGroupTiles(t, 3, Value(10))}
		result := group.IsValidSet()
		assert.True(t, result)
	})
	t.Run("ShouldReturnTrueOnRun", func(t *testing.T) {
		run := &set{tiles: createRunTiles(t, 3, 9, ColorGreen)}
		result := run.IsValidSet()
		assert.True(t, result)
	})
	t.Run("ShouldReturnFalseWithMultipleJokers", func(t *testing.T) {
		notValid := &set{tiles: []Piece{NewPiece(ValueJoker, ColorBlack), NewPiece(ValueJoker, ColorBlack), NewPiece(Value(1), ColorGreen)}}
		result := notValid.IsValidSet()
		assert.False(t, result)
	})
	t.Run("ShouldReturnFalseOnNil", func(t *testing.T) {
		result := (*set)(nil).IsValidSet()
		assert.False(t, result)
	})
	t.Run("ShouldReturnFalseOnTooFewPieces", func(t *testing.T) {
		set := &set{tiles: []Piece{NewPiece(Value(7), ColorBlack)}}
		result := set.IsValidSet()
		assert.False(t, result)
	})
	t.Run("ShouldReturnFalseOnTooManyPieces", func(t *testing.T) {
		run := &set{tiles: createRunTiles(t, 0, 18, ColorBlue)}
		result := run.IsValidSet()
		assert.False(t, result)
	})
}

func TestIsGroup(t *testing.T) {
	t.Run("ShouldReturnTrueOnGroup", func(t *testing.T) {
		group := &set{tiles: createGroupTiles(t, 3, 1)}
		result := isGroup(group)
		assert.True(t, result)
	})
	t.Run("ShouldReturnTrueOnGroupWithJoker", func(t *testing.T) {
		group := &set{tiles: []Piece{NewPiece(1, ColorBlack), NewPiece(ValueJoker, ColorBlack), NewPiece(1, ColorGreen)}}
		result := isGroup(group)
		assert.True(t, result)
	})
	t.Run("ShouldReturnFalseOnTooFewPieces", func(t *testing.T) {
		group := &set{tiles: []Piece{NewPiece(Value(1), ColorBlack), NewPiece(Value(1), ColorBlue)}}
		result := isGroup(group)
		assert.False(t, result)
	})
	t.Run("ShouldReturnFalseOnWrongColor", func(t *testing.T) {
		notGroup := &set{tiles: []Piece{NewPiece(1, ColorBlack), NewPiece(7, ColorBlack), NewPiece(3, ColorGreen)}}
		result := isGroup(notGroup)
		assert.False(t, result)
	})
}

func TestIsRun(t *testing.T) {
	t.Run("ShouldReturnTrueOnRun", func(t *testing.T) {
		run := &set{tiles: createRunTiles(t, 4, 9, ColorRed)}
		result := isRun(run)
		assert.True(t, result)
	})
	t.Run("ShouldReturnTrueOnRunWithJoker", func(t *testing.T) {
		run := &set{tiles: []Piece{NewPiece(Value(4), ColorRed), NewPiece(ValueJoker, ColorRed), NewPiece(Value(6), ColorRed)}}
		result := isRun(run)
		assert.True(t, result)
	})
	t.Run("ShouldReturnTrueOnRunStartingWithJoker", func(t *testing.T) {
		run := &set{tiles: []Piece{NewPiece(ValueJoker, ColorBlack), NewPiece(Value(4), ColorRed), NewPiece(Value(5), ColorRed)}}
		result := isRun(run)
		assert.True(t, result)
	})
	t.Run("ShouldReturnFalseOnRunWithBadJoker", func(t *testing.T) {
		run := &set{tiles: []Piece{NewPiece(Value(4), ColorRed), NewPiece(ValueJoker, ColorRed), NewPiece(Value(5), ColorRed)}}
		result := isRun(run)
		assert.False(t, result)
	})
	t.Run("ShouldReturnFalseOnTooFewPieces", func(t *testing.T) {
		run := &set{tiles: []Piece{NewPiece(Value(6), ColorBlack), NewPiece(Value(7), ColorBlack)}}
		result := isRun(run)
		assert.False(t, result)
	})
	t.Run("ShouldReturnFalseOnNotSorted", func(t *testing.T) {
		notRun := &set{tiles: []Piece{NewPiece(Value(1), ColorBlack), NewPiece(Value(2), ColorBlack), NewPiece(Value(4), ColorBlack)}}
		result := isRun(notRun)
		assert.False(t, result)
	})
	t.Run("ShouldReturnFalseOnWrongColor", func(t *testing.T) {
		notRun := &set{tiles: []Piece{NewPiece(Value(1), ColorRed), NewPiece(Value(2), ColorBlack), NewPiece(Value(3), ColorBlack)}}
		result := isRun(notRun)
		assert.False(t, result)
	})
}

func TestPiece(t *testing.T) {
	t.Run("ShouldReturnPiece", func(t *testing.T) {
		expectedPiece := NewPiece(Value(1), ColorBlack)
		get := &set{tiles: []Piece{expectedPiece}}
		result, err := get.Piece(0)
		assert.NoError(t, err)
		assert.Same(t, result.(*piece), expectedPiece)
	})
	t.Run("ShouldReturnErrorOnEmptySet", func(t *testing.T) {
		get := new(set)
		result, err := get.Piece(0)
		assert.EqualError(t, err, InvalidSet)
		assert.Nil(t, result)
	})
	t.Run("ShouldReturnErrorOnInvalidIndex", func(t *testing.T) {
		get := &set{tiles: []Piece{NewPiece(Value(1), ColorBlack)}}
		piece, err := get.Piece(-1)
		assert.EqualError(t, err, IndexOutOfBounds(-1, 1))
		assert.Nil(t, piece)
		piece, err = get.Piece(1)
		assert.EqualError(t, err, IndexOutOfBounds(-1, 1))
		assert.Nil(t, piece)
	})
}

func TestFindIndex(t *testing.T) {
	t.Run("ShouldReturnIndex", func(t *testing.T) {
		piece := NewPiece(Value(3), ColorGreen)
		search := &set{tiles: []Piece{NewPiece(Value(8), ColorBlue), NewPiece(Value(3), ColorGreen), piece}}
		index := search.findIndex(piece)
		assert.Equal(t, index, 2)
		assert.Same(t, search.tiles[index], piece)
	})
	t.Run("ShouldReturnInvalidOnNotExistingPiece", func(t *testing.T) {
		piece := NewPiece(Value(3), ColorGreen)
		search := &set{tiles: []Piece{NewPiece(Value(8), ColorBlue), NewPiece(Value(3), ColorGreen)}}
		index := search.findIndex(piece)
		assert.Negative(t, index)
	})
	t.Run("ShouldReturnInvalidOnEmptySet", func(t *testing.T) {
		piece := NewPiece(Value(3), ColorGreen)
		search := new(set)
		index := search.findIndex(piece)
		assert.Negative(t, index)
	})
}

func TestInsertPiece(t *testing.T) {
	t.Run("ShouldInsertIntoEmptySet", func(t *testing.T) {
		emptySet := new(set)
		piece := NewPiece(Value(8), ColorBlack)
		emptySet.insertPiece(piece, 0)
		assert.Len(t, emptySet.tiles, 1)
		assert.Same(t, emptySet.tiles[0], piece)
	})
	t.Run("ShouldInsertAtStart", func(t *testing.T) {
		set := &set{tiles: []Piece{NewPiece(Value(7), ColorBlue), NewPiece(Value(10), ColorGreen), NewPiece(Value(4), ColorRed)}}
		piece := NewPiece(Value(2), ColorBlack)
		start := 0
		set.insertPiece(piece, start)
		assert.Len(t, set.tiles, 4)
		assert.Same(t, set.tiles[start], piece)
	})
	t.Run("ShouldInsertAtMiddle", func(t *testing.T) {
		set := &set{tiles: []Piece{NewPiece(Value(7), ColorBlue), NewPiece(Value(10), ColorGreen), NewPiece(Value(4), ColorRed)}}
		piece := NewPiece(Value(2), ColorBlack)
		middle := 1
		set.insertPiece(piece, middle)
		assert.Len(t, set.tiles, 4)
		assert.Same(t, set.tiles[middle], piece)
	})
	t.Run("ShouldInsertAtEnd", func(t *testing.T) {
		set := &set{tiles: []Piece{NewPiece(Value(7), ColorBlue), NewPiece(Value(10), ColorGreen), NewPiece(Value(4), ColorRed)}}
		piece := NewPiece(Value(2), ColorBlack)
		end := 2
		set.insertPiece(piece, end)
		assert.Len(t, set.tiles, 4)
		assert.Same(t, set.tiles[end], piece)
	})
}

func TestInsert(t *testing.T) {
	t.Run("ShouldInsertPiece", func(t *testing.T) {
		existingPiece := NewPiece(Value(5), ColorBlack)
		original := &set{tiles: []Piece{existingPiece}}
		piece := NewPiece(Value(6), ColorBlack)
		inserted, err := original.Insert(piece, 1)
		assert.Len(t, original.tiles, 1)
		assert.NoError(t, err)
		assert.NotNil(t, inserted)
		assert.NotSame(t, inserted, original)
		assert.Len(t, inserted.(*set).tiles, 2)
		assert.Same(t, inserted.(*set).tiles[1], piece)
		assert.Same(t, inserted.(*set).tiles[0], existingPiece)
	})
	t.Run("ShouldReturnErrorOnNegativeIndex", func(t *testing.T) {
		set := new(set)
		piece := NewPiece(ValueJoker, ColorBlack)
		inserted, err := set.Insert(piece, -1)
		assert.EqualError(t, err, IndexOutOfBounds(-1, 1))
		assert.Nil(t, inserted)
	})
	t.Run("ShouldReturnErrorOnIndexOutOfBounds", func(t *testing.T) {
		set := &set{tiles: []Piece{NewPiece(Value(5), ColorBlack)}}
		piece := NewPiece(ValueJoker, ColorBlack)
		inserted, err := set.Insert(piece, 73)
		assert.EqualError(t, err, IndexOutOfBounds(-1, 2))
		assert.Nil(t, inserted)
	})
	t.Run("ShouldReturnErrorOnExistingPiece", func(t *testing.T) {
		piece := NewPiece(Value(9), ColorBlack)
		set := &set{tiles: []Piece{piece}}
		inserted, err := set.Insert(piece, 1)
		assert.EqualError(t, err, InvalidPiece)
		assert.Nil(t, inserted)
	})
}

func TestRemovePiece(t *testing.T) {
	t.Run("ShouldRemoveFromStart", func(t *testing.T) {
		remainingPiece := NewPiece(Value(8), ColorBlue)
		set := &set{tiles: []Piece{NewPiece(Value(5), ColorBlack), remainingPiece}}
		set.removePiece(0)
		assert.Len(t, set.tiles, 1)
		assert.Same(t, set.tiles[0], remainingPiece)
	})
	t.Run("ShouldRemoveFromMiddle", func(t *testing.T) {
		startPiece, endPiece := NewPiece(Value(2), ColorGreen), NewPiece(Value(8), ColorBlue)
		set := &set{tiles: []Piece{startPiece, NewPiece(Value(5), ColorBlack), endPiece}}
		set.removePiece(1)
		assert.Len(t, set.tiles, 2)
		assert.Same(t, set.tiles[0], startPiece)
		assert.Same(t, set.tiles[1], endPiece)
	})
	t.Run("ShouldRemoveFromEnd", func(t *testing.T) {
		startPiece, middlePiece := NewPiece(Value(2), ColorGreen), NewPiece(Value(8), ColorBlue)
		set := &set{tiles: []Piece{startPiece, middlePiece, NewPiece(Value(5), ColorBlack)}}
		set.removePiece(2)
		assert.Len(t, set.tiles, 2)
		assert.Same(t, set.tiles[0], startPiece)
		assert.Same(t, set.tiles[1], middlePiece)
	})
}

func TestRemove(t *testing.T) {
	t.Run("ShouldRemovePiece", func(t *testing.T) {
		piece := NewPiece(Value(12), ColorRed)
		original := &set{tiles: []Piece{piece}}
		removed, err := original.Remove(piece)
		assert.Len(t, original.tiles, 1)
		assert.NoError(t, err)
		assert.NotNil(t, removed)
		assert.Empty(t, removed.(*set).tiles)
	})
	t.Run("ShouldReturnErrorOnEmptySet", func(t *testing.T) {
		piece := NewPiece(Value(12), ColorRed)
		emptySet := new(set)
		removed, err := emptySet.Remove(piece)
		assert.Empty(t, emptySet.tiles)
		assert.EqualError(t, err, InvalidSet)
		assert.Nil(t, removed)
	})
	t.Run("ShouldReturnErrorOnNotExistingPiece", func(t *testing.T) {
		piece := NewPiece(Value(12), ColorRed)
		original := &set{tiles: []Piece{NewPiece(Value(10), ColorBlack)}}
		removed, err := original.Remove(piece)
		assert.Len(t, original.tiles, 1)
		assert.EqualError(t, err, InvalidPiece)
		assert.Nil(t, removed)
	})
}

func TestSplit(t *testing.T) {
	t.Run("ShouldSplitSet", func(t *testing.T) {
		left, right := NewPiece(Value(1), ColorBlack), NewPiece(Value(2), ColorBlack)
		original := &set{tiles: []Piece{left, right}}
		lower, upper, err := original.Split(1)
		assert.Len(t, original.tiles, 2)
		assert.NoError(t, err)
		assert.NotNil(t, lower)
		assert.NotNil(t, upper)
		assert.NotSame(t, original, lower)
		assert.NotSame(t, original, upper)
		assert.Len(t, lower.(*set).tiles, 1)
		assert.Same(t, lower.(*set).tiles[0], left)
		assert.Len(t, upper.(*set).tiles, 1)
		assert.Same(t, upper.(*set).tiles[0], right)
	})
	t.Run("ShouldReturnErrorOnTooFewPieces", func(t *testing.T) {
		original := &set{tiles: []Piece{NewPiece(Value(5), ColorGreen)}}
		lower, upper, err := original.Split(0)
		assert.Len(t, original.tiles, 1)
		assert.EqualError(t, err, TooFewPieces)
		assert.Nil(t, lower)
		assert.Nil(t, upper)
	})
	t.Run("ShouldReturnErrorOnInvalidIndex", func(t *testing.T) {
		original := &set{tiles: []Piece{NewPiece(Value(1), ColorBlack), NewPiece(Value(2), ColorBlack)}}
		lower, upper, err := original.Split(-1)
		assert.Len(t, original.tiles, 2)
		assert.EqualError(t, err, IndexOutOfBounds(0, 2))
		assert.Nil(t, lower)
		assert.Nil(t, upper)
		lower, upper, err = original.Split(2)
		assert.Len(t, original.tiles, 2)
		assert.EqualError(t, err, IndexOutOfBounds(0, 2))
		assert.Nil(t, lower)
		assert.Nil(t, upper)
	})
}

func TestCombine(t *testing.T) {
	first, second, third := NewPiece(Value(13), ColorRed), NewPiece(Value(7), ColorBlack), NewPiece(Value(1), ColorGreen)
	combined := Combine(first, second, third)
	expectedPieces := []Piece{first, second, third}
	assert.NotNil(t, combined)
	assert.IsType(t, (*set)(nil), combined)
	assert.ElementsMatch(t, combined.(*set).tiles, expectedPieces)
}
