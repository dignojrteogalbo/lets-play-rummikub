package command

import (
	"lets-play-rummikub/internal/constants"
	"lets-play-rummikub/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemove(t *testing.T) {
	t.Run("ShouldReturnRemove", func(t *testing.T) {
		game := model.NewGame(1)
		piece := model.NewPiece(model.Value(1), model.ColorBlack)
		set := model.Combine(piece, model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		setBoard(game, set)
		command, err := Remove(game, "0 0")
		assert.NoError(t, err)
		assert.NotNil(t, command)
		result := command.(*remove)
		assert.Same(t, result.game, game)
		assert.Same(t, result.set, set)
		assert.Equal(t, result.piece, piece)
	})
	t.Run("ShouldReturnErrorOnTooFewArguments", func(t *testing.T) {
		game := model.NewGame(1)
		piece := model.NewPiece(model.Value(1), model.ColorBlack)
		set := model.Combine(piece, model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		setBoard(game, set)
		command, err := Remove(game, "toofewarguments")
		assert.EqualError(t, err, constants.TooFewArguments)
		assert.Nil(t, command)
	})
	t.Run("ShouldReturnErrorOnBadSetIndex", func(t *testing.T) {
		game := model.NewGame(1)
		piece := model.NewPiece(model.Value(1), model.ColorBlack)
		set := model.Combine(piece, model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		setBoard(game, set)
		command, err := Remove(game, "bad 0")
		assert.EqualError(t, err, constants.InvalidNumberInput)
		assert.Nil(t, command)
	})
	t.Run("ShouldReturnErrorOnInvalidSet", func(t *testing.T) {
		game := model.NewGame(1)
		piece := model.NewPiece(model.Value(1), model.ColorBlack)
		set := model.Combine(piece, model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		setBoard(game, set)
		command, err := Remove(game, "1 0")
		assert.EqualError(t, err, constants.InvalidSetSelection)
		assert.Nil(t, command)
	})
	t.Run("ShouldReturnErrorOnBadPieceIndex", func(t *testing.T) {
		game := model.NewGame(1)
		piece := model.NewPiece(model.Value(1), model.ColorBlack)
		set := model.Combine(piece, model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		setBoard(game, set)
		command, err := Remove(game, "0 bad")
		assert.EqualError(t, err, constants.InvalidNumberInput)
		assert.Nil(t, command)
	})
	t.Run("ShouldReturnErrorOnInvalidPiece", func(t *testing.T) {
		game := model.NewGame(1)
		piece := model.NewPiece(model.Value(1), model.ColorBlack)
		set := model.Combine(piece, model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		setBoard(game, set)
		command, err := Remove(game, "0 5")
		assert.EqualError(t, err, constants.InvalidPieceSelection)
		assert.Nil(t, command)
	})
}

func TestInvokeRemove(t *testing.T) {
	t.Run("ShouldRemovePiece", func(t *testing.T) {
		game := model.NewGame(1)
		piece := model.NewPiece(model.Value(1), model.ColorBlack)
		set := model.Combine(piece, model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		setBoard(game, set)
		command, err := Remove(game, "0 0")
		assert.NoError(t, err)
		command.Invoke()
		gameState := unmarshal(t, game)
		assert.Len(t, gameState["board"], 1)
		assert.Len(t, gameState["board"].([]any)[0].(map[string]any)["pieces"], 3)
		assert.Len(t, gameState["piece"], 1)
	})
	t.Run("ShouldDoNothingOnBadRemove", func(t *testing.T) {
		game := model.NewGame(1)
		badPiece := model.NewPiece(model.Value(1), model.ColorBlack)
		set := model.Combine(model.NewPiece(model.Value(1), model.ColorBlack), model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		setBoard(game, set)
		command, err := Remove(game, "0 0")
		assert.NoError(t, err)
		command.(*remove).piece = badPiece
		command.Invoke()
		gameState := unmarshal(t, game)
		assert.Len(t, gameState["board"], 1)
		assert.Len(t, gameState["board"].([]any)[0].(map[string]any)["pieces"], 4)
		assert.Len(t, gameState["piece"], 0)
	})
}

func TestUndoRemove(t *testing.T) {
	game := model.NewGame(1)
	piece := model.NewPiece(model.Value(1), model.ColorBlack)
	set := model.Combine(piece, model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
	setBoard(game, set)
	command, err := Remove(game, "0 0")
	assert.NoError(t, err)
	command.Invoke()
	command.Undo()
	gameState := unmarshal(t, game)
	assert.Len(t, gameState["board"], 1)
	assert.Len(t, gameState["board"].([]any)[0].(map[string]any)["pieces"], 4)
	assert.Len(t, gameState["piece"], 0)
}
