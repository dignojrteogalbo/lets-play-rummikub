package model

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScore(t *testing.T) {
	player := &player{rack: createGroupTiles(t, 3, Value(5))}
	score := player.Score()
	assert.Equal(t, score, uint16(15))
}

func TestDealPiece(t *testing.T) {
	t.Run("ShouldAddPieceToRack", func(t *testing.T) {
		player := new(player)
		player.DealPiece(NewPiece(Value(1), ColorBlue))
		assert.Len(t, player.rack, 1)
	})
	t.Run("ShouldDoNothing", func(t *testing.T) {
		player := new(player)
		player.DealPiece(nil)
		assert.Empty(t, player.rack)
	})
}

func TestParseInt(t *testing.T) {
	t.Run("ShouldReturnInt", func(t *testing.T) {
		result, err := parseInt("2")
		assert.NoError(t, err)
		assert.Equal(t, result, 2)
	})
	t.Run("ShouldReturnError", func(t *testing.T) {
		result, err := parseInt("2\n")
		assert.EqualError(t, err, InvalidNumberInput)
		assert.Equal(t, result, -1)
	})
}

func readFromStringInput(t *testing.T, input string) {
	Reader = bufio.NewReader(strings.NewReader(input))
	t.Cleanup(func() {
		Reader = bufio.NewReader(os.Stdin)
	})
}

func readerError(t *testing.T) {
	Reader = bufio.NewReader(&strings.Reader{})
	t.Cleanup(func() {
		Reader = bufio.NewReader(os.Stdin)
	})
}

func TestSelectRackPiece(t *testing.T) {
	player := &player{rack: []Piece{NewPiece(Value(13), ColorGreen)}}
	t.Run("ShouldReturnPieceFromRack", func(t *testing.T) {
		piece, err := player.selectRackPiece("0")
		assert.NoError(t, err)
		assert.Same(t, piece, player.rack[0])
	})
	t.Run("ShouldReturnErrorOnBadInput", func(t *testing.T) {
		piece, err := player.selectRackPiece("not a number")
		assert.Error(t, err)
		assert.Nil(t, piece)
	})
	t.Run("ShouldReturnErrorOnBadIndex", func(t *testing.T) {
		piece, err := player.selectRackPiece("5")
		assert.Error(t, err)
		assert.Nil(t, piece)
	})
}

func TestSelectLoosePiece(t *testing.T) {
	game := &instance{loose: []Piece{NewPiece(Value(11), ColorBlue)}}
	t.Run("ShouldReturnLoosePiece", func(t *testing.T) {
		piece, err := selectLoosePiece("0", game)
		assert.NoError(t, err)
		assert.Len(t, game.loose, 1)
		assert.Same(t, piece, game.loose[0])
	})
	t.Run("ShouldReturnErrorOnBadInput", func(t *testing.T) {
		piece, err := selectLoosePiece("not a number", game)
		assert.Error(t, err)
		assert.Nil(t, piece)
	})
	t.Run("ShouldReturnErrorOnBadIndex", func(t *testing.T) {
		piece, err := selectLoosePiece("5", game)
		assert.Error(t, err)
		assert.Nil(t, piece)
	})
}

func TestPromptForSet(t *testing.T) {
	player := new(player)
	set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
	game := &instance{board: []Set{set}}
	t.Run("ShouldReturnSet", func(t *testing.T) {
		readFromStringInput(t, "0\n")
		selection, err := player.promptForSet(game)
		assert.NoError(t, err)
		assert.Same(t, selection, game.board[0])
	})
	t.Run("ShouldReturnErrorOnBadInput", func(t *testing.T) {
		readFromStringInput(t, "bad input\n")
		selection, err := player.promptForSet(game)
		assert.EqualError(t, err, InvalidNumberInput)
		assert.Nil(t, selection)
	})
	t.Run("ShouldReturnErrorOnInvalidSelection", func(t *testing.T) {
		readFromStringInput(t, "5\n")
		selection, err := player.promptForSet(game)
		assert.EqualError(t, err, IndexOutOfBounds(-1, 1, "set"))
		assert.Nil(t, selection)
	})
	t.Run("ShouldReturnErrorOnReaderError", func(t *testing.T) {
		readerError(t)
		selection, err := player.promptForSet(game)
		assert.Error(t, err)
		assert.Nil(t, selection)
	})
}

func TestSetPrompt(t *testing.T) {
	set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
	game := &instance{board: []Set{set}}
	player := new(player)
	t.Run("ShouldReturnSet", func(t *testing.T) {
		readFromStringInput(t, "0\n")
		selected, err := player.promptForSet(game)
		assert.NoError(t, err)
		assert.Same(t, selected, set)
	})
	t.Run("ShouldReturnErrorOnBadInput", func(t *testing.T) {
		readFromStringInput(t, "bad input\n")
		selected, err := player.promptForSet(game)
		assert.EqualError(t, err, InvalidNumberInput)
		assert.Nil(t, selected)
	})
	t.Run("ShouldReturnErrorOnSelection", func(t *testing.T) {
		readFromStringInput(t, "73\n")
		selected, err := player.promptForSet(game)
		assert.EqualError(t, err, IndexOutOfBounds(-1, 1, "set"))
		assert.Nil(t, selected)
	})
}

func TestPlayerInsert(t *testing.T) {
	t.Run("ShouldInsertIntoSet", func(t *testing.T) {
		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{board: []Set{existingSet}}
		player := &player{rack: []Piece{NewPiece(Value(4), ColorGreen)}}
		readFromStringInput(t, "0\nr0\n3\n")
		err := player.insert(game)
		assert.NoError(t, err)
		assert.Empty(t, player.rack)
		assert.Len(t, game.board, 1)
		assert.Len(t, game.board[0].(*set).tiles, existingSet.Len()+1)
		assert.NotSame(t, existingSet, game.board[0])
	})
	t.Run("ShouldReturnErrorOnBadSetSelection", func(t *testing.T) {
		set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{board: []Set{set}}
		player := &player{rack: []Piece{NewPiece(Value(5), ColorGreen)}}
		readFromStringInput(t, "2\n")
		err := player.insert(game)
		assert.Error(t, err)
		assert.Len(t, player.rack, 1)
		assert.Len(t, set.tiles, 3)
	})
	t.Run("ShouldReturnErrorOnErrorPieceSelection", func(t *testing.T) {
		set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{board: []Set{set}}
		player := &player{rack: []Piece{NewPiece(Value(5), ColorGreen)}}
		readFromStringInput(t, "0\n")
		err := player.insert(game)
		assert.Error(t, err)
		assert.Len(t, player.rack, 1)
		assert.Len(t, set.tiles, 3)
	})
	t.Run("ShouldReturnErrorOnBadPieceSelection", func(t *testing.T) {
		set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{board: []Set{set}}
		player := &player{rack: []Piece{NewPiece(Value(5), ColorGreen)}}
		readFromStringInput(t, "0\n0\n")
		err := player.insert(game)
		assert.Error(t, err)
		assert.Len(t, player.rack, 1)
		assert.Len(t, set.tiles, 3)
	})
	t.Run("ShouldReturnErrorOnErrorIndexSelection", func(t *testing.T) {
		set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{board: []Set{set}}
		player := &player{rack: []Piece{NewPiece(Value(5), ColorGreen)}}
		readFromStringInput(t, "0\nr0\n")
		err := player.insert(game)
		assert.Error(t, err)
		assert.Len(t, player.rack, 1)
		assert.Len(t, set.tiles, 3)
	})
	t.Run("ShouldReturnErrorOnBadIndexSelection", func(t *testing.T) {
		set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{board: []Set{set}}
		player := &player{rack: []Piece{NewPiece(Value(5), ColorGreen)}}
		readFromStringInput(t, "0\nr0\nbad\n")
		err := player.insert(game)
		assert.Error(t, err)
		assert.Len(t, player.rack, 1)
		assert.Len(t, set.tiles, 3)
	})
	t.Run("ShouldReturnErrorOnBadInsert", func(t *testing.T) {
		set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{board: []Set{set}}
		player := &player{rack: []Piece{NewPiece(Value(5), ColorGreen)}}
		readFromStringInput(t, "0\nr0\n-1\n")
		err := player.insert(game)
		assert.Error(t, err)
		assert.Len(t, player.rack, 1)
		assert.Len(t, set.tiles, 3)
	})
}

func TestPlayerCombine(t *testing.T) {
	t.Run("ShouldCombinePiecesFromRack", func(t *testing.T) {
		game := new(instance)
		player := &player{rack: createRunTiles(t, 1, 5, ColorBlue)}
		readFromStringInput(t, "r0\nr1\nr2\nr3\nr4\ndone\n")
		err := player.combine(game)
		assert.NoError(t, err)
		assert.Empty(t, player.rack)
		assert.Len(t, game.board, 1)
		assert.Len(t, game.board[0].(*set).tiles, 5)
	})
	t.Run("ShouldCombinePiecesFromRackAndBoard", func(t *testing.T) {
		game := &instance{loose: createRunTiles(t, 1, 4, ColorRed)}
		player := &player{rack: []Piece{NewPiece(Value(5), ColorRed), NewPiece(Value(6), ColorRed)}}
		readFromStringInput(t, "p3\nr0\nr1\ndone\n")
		err := player.combine(game)
		assert.NoError(t, err)
		assert.Empty(t, player.rack)
		assert.Len(t, game.board, 1)
		assert.Len(t, game.board[0].(*set).tiles, 3)
		assert.Len(t, game.loose, 3)
	})
	t.Run("ShouldReturnErrorOnBadPiece", func(t *testing.T) {
		game := &instance{loose: createRunTiles(t, 1, 3, ColorRed)}
		player := &player{rack: []Piece{NewPiece(Value(4), ColorRed), NewPiece(Value(5), ColorRed)}}
		readFromStringInput(t, "r0\nr1\np3\ndone\n")
		err := player.combine(game)
		assert.Error(t, err)
		assert.Len(t, player.rack, 2)
		assert.Len(t, game.loose, 3)
	})
	t.Run("ShouldReturnErrorOnReaderError", func(t *testing.T) {
		game := new(instance)
		player := new(player)
		readerError(t)
		err := player.combine(game)
		assert.Error(t, err)
	})
}

func TestPlayerRemove(t *testing.T) {
	t.Run("ShouldRemoveFromSet", func(t *testing.T) {
		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{board: []Set{existingSet}}
		player := new(player)
		readFromStringInput(t, "0\n0\n")
		err := player.remove(game)
		assert.NoError(t, err)
		assert.Len(t, game.board, 1)
		assert.NotSame(t, game.board[0], existingSet)
		assert.Len(t, game.board[0].(*set).tiles, existingSet.Len()-1)
		assert.Len(t, game.loose, 1)
		assert.Same(t, existingSet.tiles[0], game.loose[0])
	})
	t.Run("ShouldReturnErrorOnBadSetSelection", func(t *testing.T) {
		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{board: []Set{existingSet}}
		player := new(player)
		readFromStringInput(t, "1\n")
		err := player.remove(game)
		assert.Error(t, err)
		assert.Len(t, game.board, 1)
		assert.Same(t, game.board[0], existingSet)
		assert.Len(t, existingSet.tiles, 3)
	})
	t.Run("ShouldReturnErrorOnReaderError", func(t *testing.T) {
		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{board: []Set{existingSet}}
		player := new(player)
		readFromStringInput(t, "0\n")
		err := player.remove(game)
		assert.Error(t, err)
		assert.Len(t, game.board, 1)
		assert.Same(t, game.board[0], existingSet)
		assert.Len(t, existingSet.tiles, 3)
	})
	t.Run("ShouldReturnErrorOnBadNumberInput", func(t *testing.T) {
		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{board: []Set{existingSet}}
		player := new(player)
		readFromStringInput(t, "0\nbad input\n")
		err := player.remove(game)
		assert.Error(t, err)
		assert.Len(t, game.board, 1)
		assert.Same(t, game.board[0], existingSet)
		assert.Len(t, existingSet.tiles, 3)
	})
	t.Run("ShouldReturnErrorOnBadPieceSelection", func(t *testing.T) {
		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{board: []Set{existingSet}}
		player := new(player)
		readFromStringInput(t, "0\n-1\n")
		err := player.remove(game)
		assert.Error(t, err)
		assert.Len(t, game.board, 1)
		assert.Same(t, game.board[0], existingSet)
		assert.Len(t, existingSet.tiles, 3)
	})
	t.Run("ShouldReturnErrorOnBadRemove", func(t *testing.T) {
		existingSet := &set{tiles: []Piece{nil}}
		game := &instance{board: []Set{existingSet}}
		player := new(player)
		readFromStringInput(t, "0\n0\n")
		err := player.remove(game)
		assert.Error(t, err)
		assert.Len(t, game.board, 1)
		assert.Same(t, game.board[0], existingSet)
		assert.Len(t, existingSet.tiles, 1)
	})
}

func TestPlayerSplit(t *testing.T) {
	t.Run("ShouldSplitSet", func(t *testing.T) {
		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{board: []Set{existingSet}}
		player := new(player)
		readFromStringInput(t, "0\n1\n")
		err := player.split(game)
		assert.NoError(t, err)
		assert.Len(t, game.board, 2)
		assert.NotSame(t, game.board[0], existingSet)
		assert.NotSame(t, game.board[1], existingSet)
		assert.Len(t, game.board[0].(*set).tiles, 1)
		assert.Len(t, game.board[1].(*set).tiles, 2)
	})
	t.Run("ShouldReturnErrorOnBadSetSelection", func(t *testing.T) {
		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{board: []Set{existingSet}}
		player := new(player)
		readFromStringInput(t, "1\n")
		err := player.split(game)
		assert.Error(t, err)
		assert.Len(t, game.board, 1)
		assert.Same(t, game.board[0], existingSet)
		assert.Len(t, existingSet.tiles, 3)
	})
	t.Run("ShouldReturnErrorOnCannotSplitSet", func(t *testing.T) {
		existingSet := &set{tiles: []Piece{NewPiece(Value(7), ColorBlack)}}
		game := &instance{board: []Set{existingSet}}
		player := new(player)
		readFromStringInput(t, "0\n")
		err := player.split(game)
		assert.Error(t, err)
		assert.Len(t, game.board, 1)
		assert.Same(t, game.board[0], existingSet)
		assert.Len(t, existingSet.tiles, 1)
	})
	t.Run("ShouldReturnErrorOnReaderError", func(t *testing.T) {
		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{board: []Set{existingSet}}
		player := new(player)
		readFromStringInput(t, "0\n")
		err := player.split(game)
		assert.Error(t, err)
		assert.Len(t, game.board, 1)
		assert.Same(t, game.board[0], existingSet)
		assert.Len(t, existingSet.tiles, 3)
	})
	t.Run("ShouldReturnErrorOnBadIndexInput", func(t *testing.T) {
		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{board: []Set{existingSet}}
		player := new(player)
		readFromStringInput(t, "0\nbad input\n")
		err := player.split(game)
		assert.Error(t, err)
		assert.Len(t, game.board, 1)
		assert.Same(t, game.board[0], existingSet)
		assert.Len(t, existingSet.tiles, 3)
	})
	t.Run("ShouldReturnErrorOnBadSplit", func(t *testing.T) {
		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{board: []Set{existingSet}}
		player := new(player)
		readFromStringInput(t, "0\n3\n")
		err := player.split(game)
		assert.Error(t, err)
		assert.Len(t, game.board, 1)
		assert.Same(t, game.board[0], existingSet)
		assert.Len(t, existingSet.tiles, 3)
	})
}
