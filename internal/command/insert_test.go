package command

import (
	"lets-play-rummikub/internal/constants"
	"lets-play-rummikub/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	game := model.NewGame(1)
	player := game.CurrentPlayer()
	piece := model.NewPiece(model.Value(1), model.ColorBlack)
	set := model.Combine(model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
	dealPieces(player, piece)
	setBoard(game, set)
	t.Run("ShouldReturnInsert", func(t *testing.T) {
		command, err := Insert(player, game, "0 r0 0")
		assert.NoError(t, err)
		assert.NotNil(t, command)
		result := command.(*insert)
		assert.Same(t, result.game, game)
		assert.Same(t, result.player, player)
		assert.Same(t, result.piece, piece)
		assert.Same(t, result.set, set)
		assert.Equal(t, result.index, 0)
	})
	t.Run("ShouldReturnErrorOnTooFewArguments", func(t *testing.T) {
		command, err := Insert(player, game, "0")
		assert.EqualError(t, err, constants.TooFewArguments)
		assert.Nil(t, command)
	})
	t.Run("ShouldReturnErrorOnBadPosition", func(t *testing.T) {
		command, err := Insert(player, game, "0 r0 bad")
		assert.EqualError(t, err, constants.InvalidNumberInput)
		assert.Nil(t, command)
	})
	t.Run("ShouldReturnErrorOnBadSet", func(t *testing.T) {
		command, err := Insert(player, game, "bad r0 0")
		assert.EqualError(t, err, constants.InvalidNumberInput)
		assert.Nil(t, command)
	})
	t.Run("ShouldReturnErrorOnBadPieceIndex", func(t *testing.T) {
		command, err := Insert(player, game, "0 bad 0")
		assert.EqualError(t, err, constants.InvalidNumberInput)
		assert.Nil(t, command)
	})
	t.Run("ShouldReturnErrorOnBadPieceNumber", func(t *testing.T) {
		command, err := Insert(player, game, "0 bad 0")
		assert.EqualError(t, err, constants.InvalidNumberInput)
		assert.Nil(t, command)
	})
	t.Run("ShouldReturnErrorOnBadPieceSelection", func(t *testing.T) {
		command, err := Insert(player, game, "0 s0 0")
		assert.EqualError(t, err, constants.InvalidPieceSelection)
		assert.Nil(t, command)
	})
	t.Run("ShouldReturnErrorOnInvalidSet", func(t *testing.T) {
		command, err := Insert(player, game, "1 r0 0")
		assert.EqualError(t, err, constants.InvalidSetSelection)
		assert.Nil(t, command)
	})
	t.Run("ShouldReturnErrorOnInvalidPiece", func(t *testing.T) {
		command, err := Insert(player, game, "0 r1 0")
		assert.EqualError(t, err, constants.InvalidPieceSelection)
		assert.Nil(t, command)
	})
}

func TestInvokeInsert(t *testing.T) {
	t.Run("ShouldInsertRackPiece", func(t *testing.T) {
		game := model.NewGame(1)
		player := game.CurrentPlayer()
		piece := model.NewPiece(model.Value(1), model.ColorBlack)
		set := model.Combine(model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		dealPieces(player, piece)
		setBoard(game, set)
		command, err := Insert(player, game, "0 r0 0")
		assert.NoError(t, err)
		command.Invoke()
		gameState, playerState := unmarshal(t, game), unmarshal(t, player)
		assert.Len(t, gameState["board"], 1)
		assert.Len(t, gameState["board"].([]any)[0].(map[string]any)["pieces"], 4)
		assert.Len(t, gameState["piece"], 0)
		assert.Len(t, playerState["rack"], 0)
	})
	t.Run("ShouldInsertLoosePiece", func(t *testing.T) {
		game := model.NewGame(1)
		player := game.CurrentPlayer()
		piece := model.NewPiece(model.Value(1), model.ColorBlack)
		set := model.Combine(model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		setLoosePieces(game, piece)
		setBoard(game, set)
		command, err := Insert(player, game, "0 p0 0")
		assert.NoError(t, err)
		command.Invoke()
		gameState, playerState := unmarshal(t, game), unmarshal(t, player)
		assert.Len(t, gameState["board"], 1)
		assert.Len(t, gameState["board"].([]any)[0].(map[string]any)["pieces"], 4)
		assert.Len(t, gameState["piece"], 0)
		assert.Len(t, playerState["rack"], 0)
	})
	t.Run("ShouldDoNothingOnBadInsert", func(t *testing.T) {
		game := model.NewGame(1)
		player := game.CurrentPlayer()
		piece := model.NewPiece(model.Value(1), model.ColorBlack)
		set := model.Combine(model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		setLoosePieces(game, piece)
		setBoard(game, set)
		command, err := Insert(player, game, "0 p0 5")
		assert.NoError(t, err)
		command.Invoke()
		gameState, playerState := unmarshal(t, game), unmarshal(t, player)
		assert.Len(t, gameState["board"], 1)
		assert.Len(t, gameState["board"].([]any)[0].(map[string]any)["pieces"], 3)
		assert.Len(t, gameState["piece"], 1)
		assert.Len(t, playerState["rack"], 0)
	})
}

func TestUndoInsert(t *testing.T) {
	t.Run("ShouldRevertInsertRackPiece", func(t *testing.T) {
		game := model.NewGame(1)
		player := game.CurrentPlayer()
		piece := model.NewPiece(model.Value(1), model.ColorBlack)
		set := model.Combine(model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		dealPieces(player, piece)
		setBoard(game, set)
		command, err := Insert(player, game, "0 r0 0")
		assert.NoError(t, err)
		command.Invoke()
		command.Undo()
		gameState, playerState := unmarshal(t, game), unmarshal(t, player)
		assert.Len(t, gameState["board"], 1)
		assert.Len(t, gameState["board"].([]any)[0].(map[string]any)["pieces"], 3)
		assert.Len(t, gameState["piece"], 0)
		assert.Len(t, playerState["rack"], 1)
	})
	t.Run("ShouldRevertInsertLoosePiece", func(t *testing.T) {
		game := model.NewGame(1)
		player := game.CurrentPlayer()
		piece := model.NewPiece(model.Value(1), model.ColorBlack)
		set := model.Combine(model.NewPiece(model.Value(2), model.ColorBlack), model.NewPiece(model.Value(3), model.ColorBlack), model.NewPiece(model.Value(4), model.ColorBlack))
		setLoosePieces(game, piece)
		setBoard(game, set)
		command, err := Insert(player, game, "0 p0 0")
		assert.NoError(t, err)
		command.Invoke()
		command.Undo()
		gameState, playerState := unmarshal(t, game), unmarshal(t, player)
		assert.Len(t, gameState["board"], 1)
		assert.Len(t, gameState["board"].([]any)[0].(map[string]any)["pieces"], 3)
		assert.Len(t, gameState["piece"], 1)
		assert.Len(t, playerState["rack"], 0)
	})
}
