package command

import (
	"lets-play-rummikub/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombine(t *testing.T) {
	t.Run("ShouldReturnCombine", func(t *testing.T) {
		game := model.NewGame(1)
		player := game.CurrentPlayer()
		dealPieces(player, model.NewPiece(model.Value(1), model.ColorBlack), model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack))
		command, err := Combine(player, game, "r0 r1 r2")
		assert.NoError(t, err)
		assert.NotNil(t, command)
		result := command.(*combine)
		assert.Same(t, result.game, game)
		assert.Same(t, result.player, player)
		assert.Len(t, result.pieces, 3)
	})
	t.Run("ShouldReturnErrorOnBadInput", func(t *testing.T) {
		game := model.NewGame(1)
		game.DealPieces()
		player := game.CurrentPlayer()
		command, err := Combine(player, game, "badinput")
		assert.Error(t, err)
		assert.Nil(t, command)
	})
}

func TestInvokeCombine(t *testing.T) {
	t.Run("ShouldCombineLoosePieces", func(t *testing.T) {
		game := model.NewGame(1)
		player := game.CurrentPlayer()
		setLoosePieces(game, model.NewPiece(model.Value(1), model.ColorBlack), model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack))
		command, err := Combine(player, game, "p0 p1 p2")
		assert.NoError(t, err)
		command.Invoke()
		gameState, playerState := unmarshal(t, game), unmarshal(t, player)
		assert.Len(t, gameState["board"], 1)
		assert.Len(t, gameState["piece"], 0)
		assert.Len(t, playerState["rack"], 0)
	})
	t.Run("ShouldCombinePlayerRack", func(t *testing.T) {
		game := model.NewGame(1)
		player := game.CurrentPlayer()
		dealPieces(player, model.NewPiece(model.Value(1), model.ColorBlack), model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack))
		command, err := Combine(player, game, "r0 r1 r2")
		assert.NoError(t, err)
		command.Invoke()
		gameState, playerState := unmarshal(t, game), unmarshal(t, player)
		assert.Len(t, gameState["board"], 1)
		assert.Len(t, gameState["piece"], 0)
		assert.Len(t, playerState["rack"], 0)
	})
	t.Run("ShouldCombineRackAndLoosePieces", func(t *testing.T) {
		game := model.NewGame(1)
		player := game.CurrentPlayer()
		setLoosePieces(game, model.NewPiece(model.Value(1), model.ColorBlack), model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		dealPieces(player, model.NewPiece(model.Value(3), model.ColorBlack))
		command, err := Combine(player, game, "p0 p1 r0")
		assert.NoError(t, err)
		command.Invoke()
		gameState, playerState := unmarshal(t, game), unmarshal(t, player)
		assert.Len(t, gameState["board"], 1)
		assert.Len(t, gameState["piece"], 1)
		assert.Len(t, playerState["rack"], 0)
	})
}

func TestUndoCombine(t *testing.T) {
	game := model.NewGame(1)
	player := game.CurrentPlayer()
	setLoosePieces(game, model.NewPiece(model.Value(1), model.ColorBlack), model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack))
	command, err := Combine(player, game, "p0 p1 p2")
	assert.NoError(t, err)
	command.Invoke()
	command.Undo()
	gameState, playerState := unmarshal(t, game), unmarshal(t, player)
	assert.Len(t, gameState["board"], 0)
	assert.Len(t, gameState["piece"], 3)
	assert.Len(t, playerState["rack"], 0)
}
