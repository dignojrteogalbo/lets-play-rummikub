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

// func TestRemove(t *testing.T) {
// 	t.Run("ShouldReturnErrorOnTooFewPieces", func(t *testing.T) {
// 		group := &set{tiles: createGroupTiles(t, uint(3), Value(4))}
// 		piece, err := group.Piece(0)
// 		assert.NoError(t, err)
// 		err = group.Remove(piece)
// 		assert.EqualError(t, err, TooFewPieces)
// 	})
// 	t.Run("ShouldReturnErrorOnPieceNotInSet", func(t *testing.T) {
// 		group := &set{tiles: []Piece{NewPiece(4, ColorBlue), NewPiece(4, ColorBlack), NewPiece(4, ColorRed), NewPiece(4, ColorGreen)}}
// 		piece := NewPiece(4, ColorGreen)
// 		err := group.Remove(piece)
// 		assert.EqualError(t, err, InvalidPiece)
// 	})
// 	t.Run("ShouldRemoveFromGroup", func(t *testing.T) {
// 		group := &set{tiles: []Piece{NewPiece(4, ColorBlue), NewPiece(4, ColorBlack), NewPiece(4, ColorRed), NewPiece(4, ColorGreen)}}
// 		piece, err := group.Piece(2)
// 		assert.NoError(t, err)
// 		err = group.Remove(piece)
// 		assert.NoError(t, err)
// 		assert.NotContains(t, group.tiles, piece)
// 	})
// 	t.Run("ShouldReturnErrorOnInvalidSet", func(t *testing.T) {
// 		run := &set{tiles: createRunTiles(t, 4, 8, ColorRed)}
// 		piece, err := run.Piece(2)
// 		assert.NoError(t, err)
// 		err = run.Remove(piece)
// 		assert.EqualError(t, err, InvalidSet)
// 	})
// }

// func TestSplit(t *testing.T) {
// 	t.Run("ShouldReturnErrorOnSplitGroup", func(t *testing.T) {
// 		group := &set{tiles: createGroupTiles(t, 4, Value(1))}
// 		piece := NewPiece(Value(1), ColorGreen)
// 		split, err := group.Split(piece)
// 		assert.Nil(t, split)
// 		assert.EqualError(t, err, CannotSplit)
// 	})
// 	t.Run("ShouldReturnErrorOnRunLength<6", func(t *testing.T) {
// 		run := &set{tiles: createRunTiles(t, 4, 6, ColorRed)}
// 		piece := NewPiece(Value(5), ColorRed)
// 		split, err := run.Split(piece)
// 		assert.Nil(t, split)
// 		assert.EqualError(t, err, TooFewPieces)
// 	})
// 	t.Run("ShouldReturnErrorOnInvalidPiece", func(t *testing.T) {
// 		run := &set{tiles: createRunTiles(t, 4, 8, ColorBlack)}
// 		piece := NewPiece(Value(6), ColorBlue)
// 		split, err := run.Split(piece)
// 		assert.Nil(t, split)
// 		assert.EqualError(t, err, InvalidPiece)
// 	})
// 	t.Run("ShouldReturnErrorOnSmallSplit", func(t *testing.T) {
// 		run := &set{tiles: createRunTiles(t, 4, 8, ColorBlue)}
// 		piece := NewPiece(Value(5), ColorBlue)
// 		split, err := run.Split(piece)
// 		assert.Nil(t, split)
// 		assert.EqualError(t, err, TooFewPieces)
// 	})
// 	t.Run("ShouldSplitRun", func(t *testing.T) {
// 		run := &set{tiles: createRunTiles(t, 1, 9, ColorGreen)}
// 		piece := NewPiece(Value(7), ColorGreen)
// 		split, err := run.Split(piece)
// 		if assert.NoError(t, err) {
// 			assert.NotNil(t, split)
// 			assert.ElementsMatch(t, run.tiles, createRunTiles(t, 1, 7, ColorGreen))
// 			assert.ElementsMatch(t, split.(*set).tiles, createRunTiles(t, 7, 9, ColorGreen))
// 		}
// 	})
// }

// func TestCombine(t *testing.T) {
// 	t.Run("ShouldReturnErrorOnTooFewPieces", func(t *testing.T) {
// 		pieces := []Piece{NewPiece(Value(2), ColorGreen), NewPiece(Value(2), ColorBlue)}
// 		set, err := Combine(pieces...)
// 		assert.EqualError(t, err, TooFewPieces)
// 		assert.Nil(t, set)
// 	})
// 	t.Run("ShouldReturnErrorOnInvalidRun", func(t *testing.T) {
// 		pieces := []Piece{NewPiece(Value(1), ColorGreen), NewPiece(Value(2), ColorGreen), NewPiece(Value(4), ColorGreen)}
// 		invalidRun, err := Combine(pieces...)
// 		assert.EqualError(t, err, InvalidSet)
// 		assert.Nil(t, invalidRun)
// 	})
// 	t.Run("ShouldReturnErrorOnInvalidGroup", func(t *testing.T) {
// 		pieces := []Piece{NewPiece(Value(7), ColorGreen), NewPiece(Value(7), ColorBlue), NewPiece(Value(7), ColorGreen)}
// 		invalidRun, err := Combine(pieces...)
// 		assert.EqualError(t, err, InvalidSet)
// 		assert.Nil(t, invalidRun)
// 	})
// 	t.Run("ShouldReturnRun", func(t *testing.T) {
// 		pieces := createRunTiles(t, 7, 13, ColorRed)
// 		run, err := Combine(pieces...)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, run)
// 	})
// 	t.Run("ShouldReturnGroup", func(t *testing.T) {
// 		pieces := createGroupTiles(t, 4, Value(11))
// 		group, err := Combine(pieces...)
// 		assert.NoError(t, err)
// 		assert.NotNil(t, group)
// 	})
// }
