package model

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestNewPiece(t *testing.T) {
	t.Run("ShouldReturnBlackJokerPiece", func(t *testing.T) {
		newPiece := NewPiece(0, ColorBlack)
		result, ok := newPiece.(*piece)
		if assert.True(t, ok) {
			assert.Equal(t, result.value, uint8(0))
			assert.Equal(t, result.Color, ColorBlack)
		}
	})
	t.Run("ShouldReturnBlueOnePiece", func(t *testing.T) {
		newPiece := NewPiece(1, ColorBlue)
		result, ok := newPiece.(*piece)
		if assert.True(t, ok) {
			assert.Equal(t, result.value, uint8(1))
			assert.Equal(t, result.Color, ColorBlue)
		}
	})
	t.Run("ShouldReturnRedThirteenPiece", func(t *testing.T) {
		newPiece := NewPiece(13, ColorRed)
		result, ok := newPiece.(*piece)
		if assert.True(t, ok) {
			assert.Equal(t, result.value, uint8(13))
			assert.Equal(t, result.Color, ColorRed)
		}
	})
	t.Run("ShouldReturnNilOnInvalidValue", func(t *testing.T) {
		newPiece := NewPiece(14, ColorRed)
		assert.Nil(t, newPiece)
	})
	t.Run("ShouldReturnNilOnInvalidColor", func(t *testing.T) {
		newPiece := NewPiece(5, 14)
		assert.Nil(t, newPiece)
	})
}

func TestIsValidPiece(t *testing.T) {
	t.Run("ShouldReturnTrueOnValidPiece", func(t *testing.T) {
		newPiece := NewPiece(5, ColorBlue)
		valid := isValidPiece(newPiece)
		assert.True(t, valid)
	})
	t.Run("ShouldReturnFalseOnNilPiece", func(t *testing.T) {
		newPiece := (Piece)(nil)
		valid := isValidPiece(newPiece)
		assert.False(t, valid)
	})
	t.Run("ShouldReturnFalseOnFailedTypeAssertion", func(t *testing.T) {
		mockPiece := &struct{ Piece }{}
		valid := isValidPiece(mockPiece)
		assert.False(t, valid)
	})
	t.Run("ShouldReturnFalseOnInvalidValue", func(t *testing.T) {
		newPiece := &piece{14, ColorBlack}
		valid := isValidPiece(newPiece)
		assert.False(t, valid)
	})
	t.Run("ShouldReturnFalseOnInvalidColor", func(t *testing.T) {
		newPiece := &piece{0, 5}
		valid := isValidPiece(newPiece)
		assert.False(t, valid)
	})
}

func TestIsJoker(t *testing.T) {
	t.Run("ShouldReturnTrueOnJoker", func(t *testing.T) {
		joker := NewPiece(ValueJoker, ColorBlack)
		assert.True(t, joker.IsJoker())
	})
	t.Run("ShouldReturnFalseOnNotJoker", func(t *testing.T) {
		notJoker := NewPiece(11, ColorBlue)
		assert.False(t, notJoker.IsJoker())
	})
	t.Run("ShouldReturnFalseOnInvalidPiece", func(t *testing.T) {
		invalidPiece := &piece{14, ColorBlack}
		assert.False(t, invalidPiece.IsJoker())
	})
}

func TestIsSameColor(t *testing.T) {
	t.Run("ShouldReturnTrueOnSameColorPieces", func(t *testing.T) {
		blackOne := NewPiece(1, ColorBlack)
		blackThree := NewPiece(3, ColorBlack)
		assert.True(t, blackOne.IsSameColor(blackThree))
		assert.True(t, blackThree.IsSameColor(blackOne))
	})
	t.Run("ShouldReturnFalseOnDifferentColorPieces", func(t *testing.T) {
		blackOne := NewPiece(1, ColorBlack)
		blueThree := NewPiece(3, ColorBlue)
		assert.False(t, blackOne.IsSameColor(blueThree))
		assert.False(t, blueThree.IsSameColor(blackOne))
	})
	t.Run("ShouldReturnFalseOnInvalidPieces", func(t *testing.T) {
		invalidPiece := Piece(nil)
		redFive := NewPiece(5, ColorRed)
		assert.False(t, redFive.IsSameColor(invalidPiece))
	})
}

func TestIsSameValue(t *testing.T) {
	t.Run("ShouldReturnTrueOnSameValuePieces", func(t *testing.T) {
		blackOne := NewPiece(1, ColorBlack)
		redOne := NewPiece(1, ColorRed)
		assert.True(t, blackOne.IsSameValue(redOne))
		assert.True(t, redOne.IsSameValue(blackOne))
	})
	t.Run("ShouldReturnFalseOnDifferentValuePieces", func(t *testing.T) {
		blackOne := NewPiece(1, ColorBlack)
		blueThree := NewPiece(3, ColorBlue)
		assert.False(t, blackOne.IsSameValue(blueThree))
		assert.False(t, blueThree.IsSameValue(blackOne))
	})
	t.Run("ShouldReturnFalseOnInvalidPieces", func(t *testing.T) {
		invalidPiece := Piece(nil)
		redFive := NewPiece(5, ColorRed)
		assert.False(t, redFive.IsSameValue(invalidPiece))
	})
}

func TestValue(t *testing.T) {
	t.Run("ShouldReturnZeroOnInvalidPiece", func(t *testing.T) {
		piece := &piece{value: 16, Color: 5}
		value := piece.Value()
		assert.Zero(t, value)
	})
}
