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
		result := isValidSet(group)
		assert.True(t, result)
	})
	t.Run("ShouldReturnTrueOnRun", func(t *testing.T) {
		run := &set{tiles: createRunTiles(t, 3, 9, ColorGreen)}
		result := isValidSet(run)
		assert.True(t, result)
	})
	t.Run("ShouldReturnFalseOnNil", func(t *testing.T) {
		result := isValidSet(nil)
		assert.False(t, result)
	})
	t.Run("ShouldReturnFalseOnTooFewPieces", func(t *testing.T) {
		set := &set{tiles: []Piece{NewPiece(Value(7), ColorBlack)}}
		result := isValidSet(set)
		assert.False(t, result)
	})
	t.Run("ShouldReturnFalseOnTooManyPieces", func(t *testing.T) {
		run := &set{tiles: createRunTiles(t, 0, 18, ColorBlue)}
		result := isValidSet(run)
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
		run := &set{tiles: createRunTiles(t, 4, 6, ColorRed)}
		result := isRun(run)
		assert.True(t, result)
	})
	t.Run("ShouldReturnTrueOnRunWithJoker", func(t *testing.T) {
		run := &set{tiles: []Piece{NewPiece(Value(4), ColorRed), NewPiece(Value(5), ColorRed), NewPiece(Value(6), ColorRed)}}
		result := isRun(run)
		assert.True(t, result)
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

func TestInsertPiece(t *testing.T) {
	t.Run("ShouldInsertIntoEmptySet", func(t *testing.T) {
		emptySet := new(set)
		piece := NewPiece(Value(8), ColorBlack)
		emptySet.insertPiece(piece, -1)
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
		assert.EqualError(t, err, IndexOutOfBounds(len(set.tiles)))
		assert.Nil(t, inserted)
	})
	t.Run("ShouldReturnErrorOnIndexOutOfBounds", func(t *testing.T) {
		set := &set{tiles: []Piece{NewPiece(Value(5), ColorBlack)}}
		piece := NewPiece(ValueJoker, ColorBlack)
		inserted, err := set.Insert(piece, 73)
		assert.EqualError(t, err, IndexOutOfBounds(len(set.tiles)))
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
