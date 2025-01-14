package command

import (
	"lets-play-rummikub/internal/constants"
	"lets-play-rummikub/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplit(t *testing.T) {
	t.Run("ShouldReturnSplit", func(t *testing.T) {
		game := model.NewGame(1)
		set := model.Combine(model.NewPiece(model.Value(1), model.ColorBlack), model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		setBoard(game, set)
		command, err := Split(game, "0 2")
		assert.NoError(t, err)
		assert.NotNil(t, command)
		result := command.(*split)
		assert.Same(t, result.game, game)
		assert.Same(t, result.set, set)
		assert.Equal(t, result.index, 2)
	})
	t.Run("ShouldReturnErrorOnTooFewArguments", func(t *testing.T) {
		game := model.NewGame(1)
		set := model.Combine(model.NewPiece(model.Value(1), model.ColorBlack), model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		setBoard(game, set)
		command, err := Split(game, "toofewarguments")
		assert.EqualError(t, err, constants.TooFewArguments)
		assert.Nil(t, command)
	})
	t.Run("ShouldReturnErrorOnBadSetIndex", func(t *testing.T) {
		game := model.NewGame(1)
		set := model.Combine(model.NewPiece(model.Value(1), model.ColorBlack), model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		setBoard(game, set)
		command, err := Split(game, "bad 0")
		assert.EqualError(t, err, constants.InvalidNumberInput)
		assert.Nil(t, command)
	})
	t.Run("ShouldReturnErrorOnInvalidSet", func(t *testing.T) {
		game := model.NewGame(1)
		set := model.Combine(model.NewPiece(model.Value(1), model.ColorBlack), model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		setBoard(game, set)
		command, err := Split(game, "1 0")
		assert.EqualError(t, err, constants.InvalidSetSelection)
		assert.Nil(t, command)
	})
	t.Run("ShouldReturnErrorOnBadSplitIndex", func(t *testing.T) {
		game := model.NewGame(1)
		set := model.Combine(model.NewPiece(model.Value(1), model.ColorBlack), model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		setBoard(game, set)
		command, err := Split(game, "0 bad")
		assert.EqualError(t, err, constants.InvalidNumberInput)
		assert.Nil(t, command)
	})
}

func TestInvokeSplit(t *testing.T) {
	t.Run("ShouldSplitSet", func(t *testing.T) {
		game := model.NewGame(1)
		set := model.Combine(model.NewPiece(model.Value(1), model.ColorBlack), model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		setBoard(game, set)
		command, err := Split(game, "0 2")
		assert.NoError(t, err)
		assert.NotNil(t, command)
		command.Invoke()
		gameState := unmarshal(t, game)
		assert.Len(t, gameState["board"], 2)
		assert.Len(t, gameState["board"].([]any)[0].(map[string]any)["pieces"], 2)
		assert.Len(t, gameState["board"].([]any)[1].(map[string]any)["pieces"], 2)
		assert.Empty(t, gameState["piece"])
	})
	t.Run("ShouldDoNothingOnBadSplit", func(t *testing.T) {
		game := model.NewGame(1)
		set := model.Combine(model.NewPiece(model.Value(1), model.ColorBlack), model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		setBoard(game, set)
		command, err := Split(game, "0 0")
		assert.NoError(t, err)
		assert.NotNil(t, command)
		command.Invoke()
		gameState := unmarshal(t, game)
		assert.Len(t, gameState["board"], 1)
		assert.Len(t, gameState["board"].([]any)[0].(map[string]any)["pieces"], 4)
		assert.Empty(t, gameState["piece"])
	})
}

func TestUndoSplit(t *testing.T) {
	game := model.NewGame(1)
	set := model.Combine(model.NewPiece(model.Value(1), model.ColorBlack), model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
	setBoard(game, set)
	command, err := Split(game, "0 2")
	assert.NoError(t, err)
	assert.NotNil(t, command)
	command.Invoke()
	command.Undo()
	gameState := unmarshal(t, game)
	assert.Len(t, gameState["board"], 1)
	assert.Len(t, gameState["board"].([]any)[0].(map[string]any)["pieces"], 4)
	assert.Empty(t, gameState["piece"])
}
