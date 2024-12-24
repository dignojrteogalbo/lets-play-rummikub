package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func createColorTiles(color Color) []Piece {
	tiles := make([]Piece, 0)
	for value := uint8(1); value <= uint8(13); value++ {
		tiles = append(tiles, NewPiece(value, color))
		tiles = append(tiles, NewPiece(value, color))
	}
	return tiles
}

func TestNewGame(t *testing.T) {
	expectedTiles := make([]Piece, 0)
	expectedTiles = append(expectedTiles, NewPiece(ValueJoker, ColorBlack), NewPiece(ValueJoker, ColorBlack))
	expectedTiles = append(expectedTiles, createColorTiles(ColorBlack)...)
	expectedTiles = append(expectedTiles, createColorTiles(ColorBlue)...)
	expectedTiles = append(expectedTiles, createColorTiles(ColorRed)...)
	expectedTiles = append(expectedTiles, createColorTiles(ColorGreen)...)
	t.Run("ShouldCreateNewGame", func(t *testing.T) {
		game := NewGame(1)
		if assert.NotNil(t, game) {
			gameStruct, ok := game.(*instance)
			assert.True(t, ok)
			assert.ElementsMatch(t, gameStruct.tiles, expectedTiles)
		}
	})
	t.Run("ShouldReturnNilOnTooFewPlayers", func(t *testing.T) {
		game := NewGame(0)
		assert.Nil(t, game)
	})
}

func TestShuffle(t *testing.T) {
	game := NewGame(1)
	sortedTiles := make([]Piece, 0)
	sortedTiles = append(sortedTiles, NewPiece(ValueJoker, ColorBlack), NewPiece(ValueJoker, ColorBlack))
	sortedTiles = append(sortedTiles, createColorTiles(ColorBlack)...)
	sortedTiles = append(sortedTiles, createColorTiles(ColorBlue)...)
	sortedTiles = append(sortedTiles, createColorTiles(ColorRed)...)
	sortedTiles = append(sortedTiles, createColorTiles(ColorGreen)...)
	game.Shuffle()
	shuffledTiles := game.(*instance).tiles
	assert.ElementsMatch(t, shuffledTiles, sortedTiles)
	assert.NotEqual(t, shuffledTiles, sortedTiles)
}

func TestDealPieces(t *testing.T) {
	game := NewGame(2)
	game.DealPieces()
	players := game.(*instance).players
	tiles := game.(*instance).tiles
	assert.Len(t, tiles, 106-28)
	for _, p := range players {
		assert.Len(t, p.(*player).rack, 14)
	}
}
