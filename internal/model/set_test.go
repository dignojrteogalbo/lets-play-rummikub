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
		tiles = append(tiles, NewPiece(uint8(i), color))
	}
	return tiles
}

func createGroupTiles(t *testing.T, length uint, value uint8) []Piece {
	assert.True(t, length >= 3)
	tiles := make([]Piece, 0)
	for i := uint(0); i < length; i++ {
		color := Color(i%4 + 1)
		tiles = append(tiles, NewPiece(value, color))
	}
	return tiles
}

func TestIsValidSet(t *testing.T) {
	t.Run("ShouldReturnFalseOnNil", func(t *testing.T) {
		result := isValidSet(nil)
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
		run := &set{tiles: []Piece{NewPiece(4, ColorRed), NewPiece(5, ColorRed), NewPiece(6, ColorRed)}}
		result := isRun(run)
		assert.True(t, result)
	})
	t.Run("ShouldReturnFalseOnNotSorted", func(t *testing.T) {
		notRun := &set{tiles: []Piece{NewPiece(1, ColorBlack), NewPiece(2, ColorBlack), NewPiece(4, ColorBlack)}}
		result := isRun(notRun)
		assert.False(t, result)
	})
	t.Run("ShouldReturnFalseOnWrongColor", func(t *testing.T) {
		notRun := &set{tiles: []Piece{NewPiece(1, ColorRed), NewPiece(2, ColorBlack), NewPiece(3, ColorBlack)}}
		result := isRun(notRun)
		assert.False(t, result)
	})
}

func TestPiece(t *testing.T) {
	t.Run("ShouldReturnErrorOnInvalidSet", func(t *testing.T) {
		emptySet := new(set)
		result, err := emptySet.Piece(0)
		assert.EqualError(t, err, InvalidSet)
		assert.Nil(t, result)
	})
	t.Run("ShouldReturnErrorOnInvalidIndex", func(t *testing.T) {
		group := &set{createGroupTiles(t, uint(3), uint8(1))}
		result, err := group.Piece(-1)
		assert.EqualError(t, err, IndexOutOfBounds(2))
		assert.Nil(t, result)
		result, err = group.Piece(3)
		assert.EqualError(t, err, IndexOutOfBounds(2))
		assert.Nil(t, result)
	})
	t.Run("ShouldReturnPiece", func(t *testing.T) {
		run := &set{tiles: createRunTiles(t, 1, 3, ColorBlue)}
		result, err := run.Piece(0)
		assert.NoError(t, err)
		assert.Equal(t, result, NewPiece(1, ColorBlue))
		result, err = run.Piece(1)
		assert.NoError(t, err)
		assert.Equal(t, result, NewPiece(2, ColorBlue))
		result, err = run.Piece(2)
		assert.NoError(t, err)
		assert.Equal(t, result, NewPiece(3, ColorBlue))
	})
}

func TestInsert(t *testing.T) {
	// invalid cases
	t.Run("ShouldReturnErrorOnInvalidPiece", func(t *testing.T) {
		invalidPiece := &piece{value: 73, Color: 11}
		run := &set{tiles: []Piece{NewPiece(4, ColorRed), NewPiece(5, ColorRed), NewPiece(6, ColorRed)}}
		err := run.Insert(invalidPiece)
		assert.EqualError(t, err, InvalidPiece)
	})
	t.Run("ShouldReturnErrorOnInvalidSet", func(t *testing.T) {
		piece := NewPiece(3, ColorRed)
		run := &set{tiles: []Piece{NewPiece(4, ColorRed), NewPiece(5, ColorBlack), NewPiece(6, ColorRed)}}
		err := run.Insert(piece)
		assert.EqualError(t, err, InvalidSet)
	})
	// group cases
	t.Run("ShouldInsertIntoGroup", func(t *testing.T) {
		piece := NewPiece(1, ColorRed)
		group := &set{tiles: []Piece{NewPiece(1, ColorBlack), NewPiece(1, ColorBlue), NewPiece(1, ColorGreen)}}
		err := group.Insert(piece)
		assert.NoError(t, err)
		expectedTiles := []Piece{NewPiece(1, ColorBlack), NewPiece(1, ColorBlue), NewPiece(1, ColorGreen), NewPiece(1, ColorRed)}
		assert.Equal(t, group.tiles, expectedTiles)
	})
	t.Run("ShouldReturnErrorOnWrongValue", func(t *testing.T) {
		piece := NewPiece(13, ColorRed)
		group := &set{tiles: []Piece{NewPiece(1, ColorBlack), NewPiece(1, ColorBlue), NewPiece(1, ColorGreen)}}
		err := group.Insert(piece)
		assert.EqualError(t, err, WrongValueForGroup)
	})
	// run cases
	t.Run("ShouldInsertAtEndOfRun", func(t *testing.T) {
		piece := NewPiece(7, ColorRed)
		run := &set{tiles: []Piece{NewPiece(4, ColorRed), NewPiece(5, ColorRed), NewPiece(6, ColorRed)}}
		err := run.Insert(piece)
		expectedTiles := []Piece{NewPiece(4, ColorRed), NewPiece(5, ColorRed), NewPiece(6, ColorRed), NewPiece(7, ColorRed)}
		assert.NoError(t, err)
		assert.Equal(t, run.tiles, expectedTiles)
	})
	t.Run("ShouldInsertAtStartOfRun", func(t *testing.T) {
		piece := NewPiece(3, ColorRed)
		run := &set{tiles: []Piece{NewPiece(4, ColorRed), NewPiece(5, ColorRed), NewPiece(6, ColorRed)}}
		err := run.Insert(piece)
		expectedTiles := []Piece{NewPiece(3, ColorRed), NewPiece(4, ColorRed), NewPiece(5, ColorRed), NewPiece(6, ColorRed)}
		assert.NoError(t, err)
		assert.Equal(t, run.tiles, expectedTiles)
	})
	t.Run("ShouldReturnErrorOnBadInsert", func(t *testing.T) {
		piece := NewPiece(5, ColorRed)
		run := &set{tiles: []Piece{NewPiece(4, ColorRed), NewPiece(5, ColorRed), NewPiece(6, ColorRed)}}
		err := run.Insert(piece)
		assert.EqualError(t, err, CannotInsertIntoRun)
	})
	t.Run("ShouldReturnErrorOnWrongColor", func(t *testing.T) {
		piece := NewPiece(3, ColorBlack)
		run := &set{tiles: []Piece{NewPiece(4, ColorRed), NewPiece(5, ColorRed), NewPiece(6, ColorRed)}}
		err := run.Insert(piece)
		assert.EqualError(t, err, WrongColorForRun)
	})
}

func TestRemove(t *testing.T) {
	t.Run("ShouldReturnErrorOnTooFewPieces", func(t *testing.T) {
		group := &set{tiles: createGroupTiles(t, uint(3), uint8(4))}
		piece, err := group.Piece(0)
		assert.NoError(t, err)
		err = group.Remove(piece)
		assert.EqualError(t, err, TooFewPieces)
	})
	t.Run("ShouldReturnErrorOnPieceNotInSet", func(t *testing.T) {
		group := &set{tiles: []Piece{NewPiece(4, ColorBlue), NewPiece(4, ColorBlack), NewPiece(4, ColorRed), NewPiece(4, ColorGreen)}}
		piece := NewPiece(4, ColorGreen)
		err := group.Remove(piece)
		assert.EqualError(t, err, InvalidPiece)
	})
	t.Run("ShouldRemoveFromGroup", func(t *testing.T) {
		group := &set{tiles: []Piece{NewPiece(4, ColorBlue), NewPiece(4, ColorBlack), NewPiece(4, ColorRed), NewPiece(4, ColorGreen)}}
		piece, err := group.Piece(2)
		assert.NoError(t, err)
		err = group.Remove(piece)
		assert.NoError(t, err)
		assert.NotContains(t, group.tiles, piece)
	})
	t.Run("ShouldReturnErrorOnInvalidSet", func(t *testing.T) {
		run := &set{tiles: createRunTiles(t, 4, 8, ColorRed)}
		piece, err := run.Piece(2)
		assert.NoError(t, err)
		err = run.Remove(piece)
		assert.EqualError(t, err, InvalidSet)
	})
}

func TestSplit(t *testing.T) {
	t.Run("ShouldReturnErrorOnNotRun", func(t *testing.T) {
		group := &set{tiles: createGroupTiles(t, 3, 3)}
		piece := NewPiece(2, ColorBlue)
		split, err := group.Split(piece)
		assert.Nil(t, split)
		assert.EqualError(t, err, InvalidSet)
		assert.Equal(t, group.tiles, createGroupTiles(t, 3, 3))
	})
	t.Run("ShouldReturnErrorOnTooFewPieces", func(t *testing.T) {
		run := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		piece := NewPiece(4, ColorGreen)
		split, err := run.Split(piece)
		assert.Nil(t, split)
		assert.EqualError(t, err, TooFewPieces)
		assert.Equal(t, run.tiles, createRunTiles(t, 1, 3, ColorGreen))
	})
	t.Run("ShouldReturnErrorOnLateSplit", func(t *testing.T) {
		run := &set{tiles: createRunTiles(t, 1, 5, ColorGreen)}
		piece := NewPiece(4, ColorGreen)
		split, err := run.Split(piece)
		assert.Nil(t, split)
		assert.EqualError(t, err, CannotSplitRun)
		assert.Equal(t, run.tiles, createRunTiles(t, 1, 5, ColorGreen))
	})
	t.Run("ShouldReturnErrorOnEarlySplit", func(t *testing.T) {
		run := &set{tiles: createRunTiles(t, 1, 5, ColorGreen)}
		piece := NewPiece(2, ColorGreen)
		split, err := run.Split(piece)
		assert.Nil(t, split)
		assert.EqualError(t, err, CannotSplitRun)
		assert.Equal(t, run.tiles, createRunTiles(t, 1, 5, ColorGreen))
	})
	t.Run("ShouldSplitRun", func(t *testing.T) {
		piece := NewPiece(6, ColorRed)
		run := &set{tiles: createRunTiles(t, 4, 8, ColorRed)}
		split, err := run.Split(piece)
		assert.NoError(t, err)
		if assert.NotNil(t, split) {
			splitTiles := split.(*set).tiles
			expectedTiles := createRunTiles(t, 6, 8, ColorRed)
			runTiles := createRunTiles(t, 4, 6, ColorRed)
			assert.Equal(t, splitTiles, expectedTiles)
			assert.Equal(t, run.tiles, runTiles)
		}
	})
}

func TestCombine(t *testing.T) {
	// bad input cases
	t.Run("ShouldReturnErrorOnUnpairedArguments", func(t *testing.T) {
		set := &set{tiles: createGroupTiles(t, 3, 7)}
		combine, err := Combine(set)
		assert.Nil(t, combine)
		assert.EqualError(t, err, InvalidCombineArguments)
	})
	t.Run("ShouldReturnErrorOnBadArguments", func(t *testing.T) {
		set := &set{tiles: createGroupTiles(t, 3, 7)}
		combine, err := Combine(set, "hello world")
		assert.Nil(t, combine)
		assert.EqualError(t, err, InvalidCombineArguments)
		combine, err = Combine("hello world", 1)
		assert.Nil(t, combine)
		assert.EqualError(t, err, InvalidCombineArguments)
	})
	t.Run("ShouldReturnErrorOnTooFewPieces", func(t *testing.T) {
		set := &set{tiles: createGroupTiles(t, 3, 7)}
		combine, err := Combine(set, 1)
		assert.Nil(t, combine)
		assert.EqualError(t, err, TooFewPieces)
	})
	// good input cases
	t.Run("ShouldCombineSets", func(t *testing.T) {
		run := &set{tiles: createRunTiles(t, 1, 4, ColorBlack)}
		group := &set{tiles: createGroupTiles(t, 4, 1)}
		piece := NewPiece(1, ColorBlue)
		combine, err := Combine(run, 0, group, 2, piece)
		assert.NoError(t, err)
		if assert.NotNil(t, combine) {
			runTiles := createRunTiles(t, 2, 4, ColorBlack)
			groupTiles := []Piece{NewPiece(1, ColorBlack), NewPiece(1, ColorBlue), NewPiece(1, ColorGreen)}
			expectedTiles := []Piece{NewPiece(1, ColorBlack), NewPiece(1, ColorRed), NewPiece(1, ColorBlue)}
			assert.Equal(t, run.tiles, runTiles)
			assert.Equal(t, group.tiles, groupTiles)
			assert.Equal(t, combine.(*set).tiles, expectedTiles)
		}
	})
	t.Run("ShouldRevertOnBadCombination", func(t *testing.T) {
		run := &set{tiles: createRunTiles(t, 1, 4, ColorBlack)}
		group := &set{tiles: createGroupTiles(t, 4, 1)}
		revert := &set{tiles: createRunTiles(t, 1, 3, ColorBlue)}
		combine, err := Combine(run, 0, group, 2, NewPiece(1, ColorGreen), revert, 0)
		assert.EqualError(t, err, InvalidSet)
		assert.Nil(t, combine)
		assert.Equal(t, run.tiles, createRunTiles(t, 1, 4, ColorBlack))
		assert.Equal(t, group.tiles, createGroupTiles(t, 4, 1))
		assert.Equal(t, revert.tiles, createRunTiles(t, 1, 3, ColorBlue))
	})
	t.Run("ShouldReturnErrorOnInvalidCombination", func(t *testing.T) {
		pieces := []any{NewPiece(1, ColorRed), NewPiece(2, ColorGreen), NewPiece(3, ColorRed)}
		combine, err := Combine(pieces...)
		assert.EqualError(t, err, InvalidSet)
		assert.Nil(t, combine)
	})
}
