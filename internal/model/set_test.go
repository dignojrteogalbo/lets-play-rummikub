package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsGroup(t *testing.T) {
	t.Run("ShouldReturnTrueOnGroup", func(t *testing.T) {
		group := &set{tiles: []Piece{NewPiece(1, ColorBlack), NewPiece(1, ColorBlue), NewPiece(1, ColorGreen)}}
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
		run := &set{tiles: []Piece{NewPiece(4, ColorRed), NewPiece(5, ColorRed), NewPiece(6, ColorRed)}}
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
