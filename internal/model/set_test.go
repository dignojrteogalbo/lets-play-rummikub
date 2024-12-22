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
		result := group.IsGroup()
		assert.True(t, result)
	})
	t.Run("ShouldReturnTrueOnGroupWithJoker", func(t *testing.T) {
		group := &set{tiles: []Piece{NewPiece(1, ColorBlack), NewPiece(ValueJoker, ColorBlack), NewPiece(1, ColorGreen)}}
		result := group.IsGroup()
		assert.True(t, result)
	})
	t.Run("ShouldReturnFalseOnWrongColor", func(t *testing.T) {
		notGroup := &set{tiles: []Piece{NewPiece(1, ColorBlack), NewPiece(7, ColorBlack), NewPiece(3, ColorGreen)}}
		result := notGroup.IsGroup()
		assert.False(t, result)
	})
}

func TestIsRun(t *testing.T) {
	t.Run("ShouldReturnTrueOnRun", func(t *testing.T) {
		run := &set{tiles: createRunTiles(t, 4, 6, ColorRed)}
		result := run.IsRun()
		assert.True(t, result)
	})
	t.Run("ShouldReturnTrueOnRunWithJoker", func(t *testing.T) {
		run := &set{tiles: []Piece{NewPiece(4, ColorRed), NewPiece(5, ColorRed), NewPiece(6, ColorRed)}}
		result := run.IsRun()
		assert.True(t, result)
	})
	t.Run("ShouldReturnFalseOnNotSorted", func(t *testing.T) {
		notRun := &set{tiles: []Piece{NewPiece(1, ColorBlack), NewPiece(2, ColorBlack), NewPiece(4, ColorBlack)}}
		result := notRun.IsRun()
		assert.False(t, result)
	})
	t.Run("ShouldReturnFalseOnWrongColor", func(t *testing.T) {
		notRun := &set{tiles: []Piece{NewPiece(1, ColorRed), NewPiece(2, ColorBlack), NewPiece(3, ColorBlack)}}
		result := notRun.IsRun()
		assert.False(t, result)
	})
}

func TestInsert(t *testing.T) {
	// invalid cases
	t.Run("ShouldReturnErrorOnInvalidPiece", func(t *testing.T) {
		invalidPiece := &piece{value: 73, Color: 11}
		run := &set{tiles: []Piece{NewPiece(4, ColorRed), NewPiece(5, ColorRed), NewPiece(6, ColorRed)}}
		err := run.Insert(invalidPiece, 0)
		assert.EqualError(t, err, "piece is invalid")
	})
	t.Run("ShouldReturnErrorOnInvalidIndex", func(t *testing.T) {
		piece := NewPiece(3, ColorRed)
		run := &set{tiles: []Piece{NewPiece(4, ColorRed), NewPiece(5, ColorRed), NewPiece(6, ColorRed)}}
		err := run.Insert(piece, -1)
		assert.EqualError(t, err, "index must be >= 0 and <= 3")
	})
	t.Run("ShouldReturnErrorOnInvalidSet", func(t *testing.T) {
		piece := NewPiece(3, ColorRed)
		run := &set{tiles: []Piece{NewPiece(4, ColorRed), NewPiece(5, ColorBlack), NewPiece(6, ColorRed)}}
		err := run.Insert(piece, 0)
		assert.EqualError(t, err, "set is invalid")
	})
	// group cases
	t.Run("ShouldInsertAtEndOfGroup", func(t *testing.T) {
		piece := NewPiece(1, ColorRed)
		group := &set{tiles: []Piece{NewPiece(1, ColorBlack), NewPiece(1, ColorBlue), NewPiece(1, ColorGreen)}}
		err := group.Insert(piece, 3)
		assert.NoError(t, err)
		expectedTiles := []Piece{NewPiece(1, ColorBlack), NewPiece(1, ColorBlue), NewPiece(1, ColorGreen), NewPiece(1, ColorRed)}
		assert.Equal(t, group.tiles, expectedTiles)
	})
	t.Run("ShouldInsertAtMiddleOfGroup", func(t *testing.T) {
		piece := NewPiece(1, ColorRed)
		group := &set{tiles: []Piece{NewPiece(1, ColorBlack), NewPiece(1, ColorBlue), NewPiece(1, ColorGreen)}}
		err := group.Insert(piece, 1)
		assert.NoError(t, err)
		expectedTiles := []Piece{NewPiece(1, ColorBlack), NewPiece(1, ColorRed), NewPiece(1, ColorBlue), NewPiece(1, ColorGreen)}
		assert.Equal(t, group.tiles, expectedTiles)
	})
	t.Run("ShouldInsertAtStartOfGroup", func(t *testing.T) {
		piece := NewPiece(1, ColorRed)
		group := &set{tiles: []Piece{NewPiece(1, ColorBlack), NewPiece(1, ColorBlue), NewPiece(1, ColorGreen)}}
		err := group.Insert(piece, 0)
		assert.NoError(t, err)
		expectedTiles := []Piece{NewPiece(1, ColorRed), NewPiece(1, ColorBlack), NewPiece(1, ColorBlue), NewPiece(1, ColorGreen)}
		assert.Equal(t, group.tiles, expectedTiles)
	})
	t.Run("ShouldReturnErrorOnWrongValue", func(t *testing.T) {
		piece := NewPiece(13, ColorRed)
		group := &set{tiles: []Piece{NewPiece(1, ColorBlack), NewPiece(1, ColorBlue), NewPiece(1, ColorGreen)}}
		err := group.Insert(piece, 0)
		assert.EqualError(t, err, "piece does not match the value of the group")
	})
	// run cases
	t.Run("ShouldInsertAtEndOfRun", func(t *testing.T) {
		piece := NewPiece(7, ColorRed)
		run := &set{tiles: []Piece{NewPiece(4, ColorRed), NewPiece(5, ColorRed), NewPiece(6, ColorRed)}}
		err := run.Insert(piece, 3)
		expectedTiles := []Piece{NewPiece(4, ColorRed), NewPiece(5, ColorRed), NewPiece(6, ColorRed), NewPiece(7, ColorRed)}
		assert.NoError(t, err)
		assert.Equal(t, run.tiles, expectedTiles)
	})
	t.Run("ShouldReturnErrorOnMiddleInsert", func(t *testing.T) {
		piece := NewPiece(5, ColorRed)
		run := &set{tiles: []Piece{NewPiece(4, ColorRed), NewPiece(5, ColorRed), NewPiece(6, ColorRed)}}
		err := run.Insert(piece, 1)
		assert.EqualError(t, err, "piece cannot be inserted into run")
	})
	t.Run("ShouldInsertAtStartOfRun", func(t *testing.T) {
		piece := NewPiece(3, ColorRed)
		run := &set{tiles: []Piece{NewPiece(4, ColorRed), NewPiece(5, ColorRed), NewPiece(6, ColorRed)}}
		err := run.Insert(piece, 0)
		expectedTiles := []Piece{NewPiece(3, ColorRed), NewPiece(4, ColorRed), NewPiece(5, ColorRed), NewPiece(6, ColorRed)}
		assert.NoError(t, err)
		assert.Equal(t, run.tiles, expectedTiles)
	})
	t.Run("ShouldReturnErrorOnWrongColor", func(t *testing.T) {
		piece := NewPiece(3, ColorBlack)
		run := &set{tiles: []Piece{NewPiece(4, ColorRed), NewPiece(5, ColorRed), NewPiece(6, ColorRed)}}
		err := run.Insert(piece, 0)
		assert.EqualError(t, err, "piece does not match the color of the run")
	})
}

func TestRemove(t *testing.T) {
	// invalid cases
	t.Run("ShouldReturnErrorOnSetOfSize3", func(t *testing.T) {
		group := &set{tiles: []Piece{NewPiece(4, ColorBlue), NewPiece(4, ColorBlack), NewPiece(4, ColorRed)}}
		rack := []Piece{NewPiece(5, ColorBlue), NewPiece(6, ColorBlue)}
		set, err := group.Remove(0, 1, rack...)
		assert.Nil(t, set)
		assert.EqualError(t, err, TooFewPieces)
	})
	t.Run("ShouldReturnErrorOnTooFewPieces", func(t *testing.T) {
		group := &set{tiles: []Piece{NewPiece(4, ColorBlue), NewPiece(4, ColorBlack), NewPiece(4, ColorRed), NewPiece(4, ColorGreen)}}
		rack := []Piece{NewPiece(5, ColorBlue)}
		set, err := group.Remove(0, 1, rack...)
		assert.Nil(t, set)
		assert.EqualError(t, err, TooFewPieces)
	})
	t.Run("ShouldReturnErrorOnIndexOutOfBounds", func(t *testing.T) {
		group := &set{tiles: []Piece{NewPiece(4, ColorBlue), NewPiece(4, ColorBlack), NewPiece(4, ColorRed), NewPiece(4, ColorGreen)}}
		rack := []Piece{NewPiece(3, ColorBlue), NewPiece(5, ColorBlue)}
		set, err := group.Remove(4, 1, rack...)
		assert.Nil(t, set)
		assert.EqualError(t, err, IndexOutOfBounds(4))
	})
	t.Run("ShouldReturnErrorOnPostiionOutOfBounds", func(t *testing.T) {
		group := &set{tiles: []Piece{NewPiece(4, ColorBlue), NewPiece(4, ColorBlack), NewPiece(4, ColorRed), NewPiece(4, ColorGreen)}}
		rack := []Piece{NewPiece(3, ColorBlue), NewPiece(5, ColorBlue)}
		set, err := group.Remove(0, -1, rack...)
		assert.Nil(t, set)
		assert.EqualError(t, err, IndexOutOfBounds(2, "position"))
	})
	t.Run("ShouldReturnErrorOnInvalidSet", func(t *testing.T) {
		group := &set{tiles: []Piece{NewPiece(4, ColorBlue), NewPiece(4, ColorBlack), NewPiece(4, ColorRed), NewPiece(4, ColorGreen)}}
		rack := []Piece{NewPiece(3, ColorBlue), NewPiece(5, ColorBlue)}
		set, err := group.Remove(0, 0, rack...)
		assert.Nil(t, set)
		assert.EqualError(t, err, InvalidSet)
	})
	// create new set cases
	t.Run("ShouldCreateRunFromGroup", func(t *testing.T) {
		group := &set{tiles: []Piece{NewPiece(4, ColorBlue), NewPiece(4, ColorBlack), NewPiece(4, ColorRed), NewPiece(4, ColorGreen)}}
		rack := []Piece{NewPiece(3, ColorBlue), NewPiece(5, ColorBlue)}
		run, err := group.Remove(0, 1, rack...)
		assert.NoError(t, err)
		if assert.True(t, isValidSet(run)) {
			runTiles := run.(*set).tiles
			expectedTiles := []Piece{NewPiece(3, ColorBlue), NewPiece(4, ColorBlue), NewPiece(5, ColorBlue)}
			removedTiles := []Piece{NewPiece(4, ColorBlack), NewPiece(4, ColorRed), NewPiece(4, ColorGreen)}
			assert.Equal(t, runTiles, expectedTiles)
			assert.Equal(t, group.tiles, removedTiles)
		}
	})
	t.Run("ShouldCreateGroupFromRun", func(t *testing.T) {
		run := &set{tiles: []Piece{NewPiece(3, ColorBlue), NewPiece(4, ColorBlue), NewPiece(5, ColorBlue), NewPiece(6, ColorBlue)}}
		rack := []Piece{NewPiece(6, ColorBlack), NewPiece(6, ColorGreen)}
		group, err := run.Remove(3, 0, rack...)
		assert.NoError(t, err)
		if assert.True(t, isValidSet(group)) {
			groupTiles := group.(*set).tiles
			expectedTiles := []Piece{NewPiece(6, ColorBlue), NewPiece(6, ColorBlack), NewPiece(6, ColorGreen)}
			removedTiles := []Piece{NewPiece(3, ColorBlue), NewPiece(4, ColorBlue), NewPiece(5, ColorBlue)}
			assert.Equal(t, groupTiles, expectedTiles)
			assert.Equal(t, run.tiles, removedTiles)
		}
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
	t.Run("ShouldReturnErrorOnInvalidSplit", func(t *testing.T) {
		run := &set{tiles: createRunTiles(t, 1, 5, ColorGreen)}
		piece := NewPiece(7, ColorGreen)
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
