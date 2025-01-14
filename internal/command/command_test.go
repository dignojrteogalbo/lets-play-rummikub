package command

import (
	"bytes"
	"encoding/json"
	"lets-play-rummikub/internal/constants"
	"lets-play-rummikub/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setLoosePieces(game model.Game, pieces ...model.Piece) {
	for _, piece := range pieces {
		game.AddLoosePiece(piece)
	}
}

func setBoard(game model.Game, sets ...model.Set) {
	for _, set := range sets {
		game.AddSet(set)
	}
}

func dealPieces(player model.Player, pieces ...model.Piece) {
	for _, piece := range pieces {
		player.DealPiece(piece)
	}
}

func unmarshal(t *testing.T, anything json.Marshaler) map[string]any {
	jsonBytes, err := anything.MarshalJSON()
	assert.NoError(t, err)
	var result map[string]any
	err = json.NewDecoder(bytes.NewBuffer(jsonBytes)).Decode(&result)
	assert.NoError(t, err)
	return result
}

func TestParseInt(t *testing.T) {
	t.Run("ShouldReturnInt", func(t *testing.T) {
		result, err := parseInt("2")
		assert.NoError(t, err)
		assert.Equal(t, result, 2)
	})
	t.Run("ShouldReturnError", func(t *testing.T) {
		result, err := parseInt("2\n")
		assert.EqualError(t, err, constants.InvalidNumberInput)
		assert.Equal(t, result, -1)
	})
}

type mockHasPiece struct{}

func (*mockHasPiece) Piece(int) (model.Piece, error) {
	return nil, nil
}

func TestGetPieceFrom(t *testing.T) {
	game, player, set := model.NewGame(1), model.NewPlayer(), model.Combine()
	options := []model.HasPiece{game, player, set, new(mockHasPiece)}
	t.Run("ShouldGetFromPlayer", func(t *testing.T) {
		selectPiece, err := getPieceFrom('r', options...)
		assert.NoError(t, err)
		assert.Same(t, selectPiece, player)
	})
	t.Run("ShouldGetFromGame", func(t *testing.T) {
		selectPiece, err := getPieceFrom('p', options...)
		assert.NoError(t, err)
		assert.Same(t, selectPiece, game)
	})
	t.Run("ShouldGetFromSet", func(t *testing.T) {
		selectPiece, err := getPieceFrom('s', options...)
		assert.NoError(t, err)
		assert.Same(t, selectPiece, set)
	})
	t.Run("ShouldReturnErrorOnInvalidSelection", func(t *testing.T) {
		selectPiece, err := getPieceFrom('x', options...)
		assert.EqualError(t, err, constants.InvalidPieceSelection)
		assert.Nil(t, selectPiece)
	})
}

func TestParseSelectedPieces(t *testing.T) {
	a, b, c := model.NewPiece(model.Value(1), model.ColorBlack), model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack)
	player, game, set := model.NewPlayer(), model.NewGame(1), model.Combine(c)
	options := []model.HasPiece{player, game, set}
	player.DealPiece(a)
	game.AddLoosePiece(b)
	t.Run("ShouldReturnPieces", func(t *testing.T) {
		pieces, err := parseSelectedPieces("r0 p0 s0", options...)
		assert.NoError(t, err)
		assert.Equal(t, pieces, []model.Piece{a, b, c})
	})
	t.Run("ShouldReturnErrorOnBadSelector", func(t *testing.T) {
		pieces, err := parseSelectedPieces("r0 x0 s0", options...)
		assert.EqualError(t, err, constants.InvalidPieceSelection)
		assert.Nil(t, pieces)
	})
	t.Run("ShouldReturnErrorOnBadIndex", func(t *testing.T) {
		pieces, err := parseSelectedPieces("r0 p0 sinvalid", options...)
		assert.EqualError(t, err, constants.InvalidNumberInput)
		assert.Nil(t, pieces)
	})
	t.Run("ShouldReturnErrorOnInvalidPieceIndex", func(t *testing.T) {
		pieces, err := parseSelectedPieces("r2 p0 s0", options...)
		assert.EqualError(t, err, constants.InvalidPieceSelection)
		assert.Nil(t, pieces)
	})
}
