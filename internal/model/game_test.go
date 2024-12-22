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
	game := NewGame()
	if assert.NotNil(t, game) {
		gameStruct, ok := game.(*gameInstance)
		assert.True(t, ok)
		assert.ElementsMatch(t, gameStruct.tiles, expectedTiles)
	}
}
