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

func TestTotalMoves(t *testing.T) {
	player := &player{totalMoves: 73}
	assert.Equal(t, player.TotalMoves(), uint16(73))
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

func TestSelectPiece(t *testing.T) {
	game := new(instance)
	player := &player{rack: []Piece{NewPiece(Value(7), ColorGreen), NewPiece(Value(3), ColorBlue)}}
	t.Run("ShouldSelectPieceFromRack", func(t *testing.T) {
		game := &instance{sets: []Set{&set{createGroupTiles(t, 4, Value(11))}}}
		selected, err := player.selectPiece("r0\n", game)
		assert.NoError(t, err)
		assert.NotNil(t, selected)
		assert.Same(t, selected.(*piece), player.rack[0].(*piece))
	})
	t.Run("ShouldSelectPieceFromSet", func(t *testing.T) {
		game := &instance{sets: []Set{&set{createGroupTiles(t, 4, Value(11))}}}
		expectedPiece := game.sets[0].(*set).tiles[2]
		selected, err := player.selectPiece("s0,2\n", game)
		assert.NoError(t, err)
		assert.NotNil(t, selected)
		assert.Same(t, selected.(*piece), expectedPiece)
		assert.Len(t, game.sets[0].(*set).tiles, 3)
	})
	t.Run("ShouldReturnErrorOnInvalidRackInput", func(t *testing.T) {
		selected, err := player.selectPiece("rbad input\n", game)
		assert.Nil(t, selected)
		assert.Error(t, err)
	})
	t.Run("ShouldReturnErrorOnInvalidRackIndex", func(t *testing.T) {
		selected, err := player.selectPiece("r73\n", game)
		assert.Nil(t, selected)
		assert.Error(t, err)
	})
	t.Run("ShouldReturnErrorOnInvalidSelection", func(t *testing.T) {
		selected, err := player.selectPiece("s02\n", game)
		assert.Nil(t, selected)
		assert.Error(t, err)
	})
	t.Run("ShouldReturnErrorOnBadSetInput", func(t *testing.T) {
		selected, err := player.selectPiece("sbad set,2\n", game)
		assert.Nil(t, selected)
		assert.Error(t, err)
	})
	t.Run("ShouldReturnErrorOnBadPieceInput", func(t *testing.T) {
		selected, err := player.selectPiece("s0,bad piece\n", game)
		assert.Nil(t, selected)
		assert.Error(t, err)
	})
	t.Run("ShouldReturnErrorOnInvalidSetSelection", func(t *testing.T) {
		selected, err := player.selectPiece("s2,2\n", game)
		assert.Nil(t, selected)
		assert.Error(t, err)
	})
	t.Run("ShouldReturnErrorOnInvalidPieceSelection", func(t *testing.T) {
		game := &instance{sets: []Set{&set{createGroupTiles(t, 4, Value(11))}}}
		selected, err := player.selectPiece("s0,7\n", game)
		assert.Nil(t, selected)
		assert.Error(t, err)
	})
	t.Run("ShouldReturnErrorOnBadInput", func(t *testing.T) {
		selected, err := player.selectPiece("bad input\n", game)
		assert.Nil(t, selected)
		assert.Error(t, err)
	})
	t.Run("ShouldReturnErrorOnBadRemove", func(t *testing.T) {
		game := &instance{sets: []Set{&set{createRunTiles(t, 1, 3, ColorGreen)}}}
		selected, err := player.selectPiece("s0,1", game)
		assert.Nil(t, selected)
		assert.Error(t, err)
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

func TestPromptForPiece(t *testing.T) {
	player := &player{rack: []Piece{NewPiece(Value(4), ColorGreen)}}
	set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
	t.Run("ShouldReturnPieceFromRack", func(t *testing.T) {
		readFromStringInput(t, "0\n")
		piece, err := player.promptForPiece(nil)
		assert.NoError(t, err)
		assert.Same(t, piece, player.rack[0])
	})
	t.Run("ShouldReturnPieceFromSet", func(t *testing.T) {
		readFromStringInput(t, "0\n")
		piece, err := player.promptForPiece(set)
		assert.NoError(t, err)
		assert.Same(t, piece, set.tiles[0])
	})
	t.Run("ShouldReturnErrorOnBadInput", func(t *testing.T) {
		readFromStringInput(t, "bad input\n")
		piece, err := player.promptForPiece(nil)
		assert.EqualError(t, err, InvalidNumberInput)
		assert.Nil(t, piece)
	})
	t.Run("ShouldReturnErrorOnInvalidSelection", func(t *testing.T) {
		readFromStringInput(t, "5\n")
		piece, err := player.promptForPiece(nil)
		assert.EqualError(t, err, IndexOutOfBounds(len(player.rack)-1, "piece"))
		assert.Nil(t, piece)
	})
	t.Run("ShouldReturnErrorOnReaderError", func(t *testing.T) {
		readerError(t)
		piece, err := player.promptForPiece(nil)
		assert.Error(t, err)
		assert.Nil(t, piece)
	})
}

func TestPromptForSet(t *testing.T) {
	player := new(player)
	set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
	game := &instance{sets: []Set{set}}
	t.Run("ShouldReturnSet", func(t *testing.T) {
		readFromStringInput(t, "0\n")
		selection, err := player.promptForSet(game)
		assert.NoError(t, err)
		assert.Same(t, selection, game.sets[0])
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
		assert.EqualError(t, err, IndexOutOfBounds(len(game.sets)-1, "set"))
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
	game := &instance{sets: []Set{set}}
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
		assert.EqualError(t, err, IndexOutOfBounds(len(game.sets)-1, "set"))
		assert.Nil(t, selected)
	})
}

func TestStartTurn(t *testing.T) {
	t.Run("ShouldAddPieceToRackOnNoSuccessfulMoves", func(t *testing.T) {
		player := &player{rack: []Piece{NewPiece(Value(4), ColorBlack)}}
		game := &instance{tiles: []Piece{NewPiece(Value(5), ColorBlue)}}
		readFromStringInput(t, "help\nnot a command\ndone\n")
		player.StartTurn(game)
		assert.Equal(t, player.TotalMoves(), uint16(0))
		assert.Len(t, player.rack, 2)
		assert.Len(t, game.tiles, 0)
	})
	t.Run("ShouldNotAddPieceToRackOnInsert", func(t *testing.T) {
		player := &player{rack: []Piece{NewPiece(Value(11), ColorGreen)}}
		group := &set{tiles: createGroupTiles(t, 3, Value(11))}
		game := &instance{tiles: []Piece{NewPiece(Value(5), ColorBlue)}, sets: []Set{group}}
		readFromStringInput(t, "insert\n0\n0\ndone\n")
		player.StartTurn(game)
		assert.Equal(t, player.TotalMoves(), uint16(1))
		assert.Empty(t, player.rack)
		assert.Len(t, game.tiles, 1)
	})
	t.Run("ShouldNotAddPieceToRackOnCombine", func(t *testing.T) {
		player := &player{rack: createRunTiles(t, 11, 13, ColorRed)}
		game := &instance{tiles: []Piece{NewPiece(Value(5), ColorBlue)}, sets: []Set{}}
		readFromStringInput(t, "combine\nr0\nr1\nr2\ndone\ndone\n")
		player.StartTurn(game)
		assert.Equal(t, player.TotalMoves(), uint16(1))
		assert.Empty(t, player.rack)
		assert.Len(t, game.sets, 1)
	})
}

func TestPlayerInsert(t *testing.T) {
	t.Run("ShouldInsertIntoSet", func(t *testing.T) {
		set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{sets: []Set{set}}
		player := &player{rack: []Piece{NewPiece(Value(4), ColorGreen)}}
		readFromStringInput(t, "0\n0\n")
		err := player.insert(game)
		assert.NoError(t, err)
		assert.Empty(t, player.rack)
		assert.Len(t, set.tiles, 4)
	})
	t.Run("ShouldReturnErrorOnBadSetSelection", func(t *testing.T) {
		set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{sets: []Set{set}}
		player := &player{rack: []Piece{NewPiece(Value(5), ColorGreen)}}
		readFromStringInput(t, "2\n")
		err := player.insert(game)
		assert.Error(t, err)
		assert.Len(t, player.rack, 1)
		assert.Len(t, set.tiles, 3)
	})
	t.Run("ShouldReturnErrorOnBadPieceSelection", func(t *testing.T) {
		set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{sets: []Set{set}}
		player := &player{rack: []Piece{NewPiece(Value(5), ColorGreen)}}
		readFromStringInput(t, "0\n4\n")
		err := player.insert(game)
		assert.Error(t, err)
		assert.Len(t, player.rack, 1)
		assert.Len(t, set.tiles, 3)
	})
	t.Run("ShouldReturnErrorOnBadInsert", func(t *testing.T) {
		set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
		game := &instance{sets: []Set{set}}
		player := &player{rack: []Piece{NewPiece(Value(5), ColorGreen)}}
		readFromStringInput(t, "0\n0\n")
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
		assert.Len(t, game.sets, 1)
		assert.Len(t, game.sets[0].(*set).tiles, 5)
	})
	t.Run("ShouldCombinePiecesFromRackAndBoard", func(t *testing.T) {
		run := &set{tiles: createRunTiles(t, 1, 4, ColorRed)}
		game := &instance{sets: []Set{run}}
		player := &player{rack: []Piece{NewPiece(Value(5), ColorRed), NewPiece(Value(6), ColorRed)}}
		readFromStringInput(t, "s0,3\nr0\nr1\ndone\n")
		err := player.combine(game)
		assert.NoError(t, err)
		assert.Empty(t, player.rack)
		assert.Len(t, game.sets, 2)
		assert.Len(t, game.sets[0].(*set).tiles, 3)
		assert.Len(t, game.sets[1].(*set).tiles, 3)
	})
	t.Run("ShouldReturnErrorOnBadPiece", func(t *testing.T) {
		run := &set{tiles: createRunTiles(t, 1, 3, ColorRed)}
		game := &instance{sets: []Set{run}}
		player := &player{rack: []Piece{NewPiece(Value(4), ColorRed), NewPiece(Value(5), ColorRed)}}
		readFromStringInput(t, "s0,2\nr0\nr1\ndone\n")
		err := player.combine(game)
		assert.Error(t, err)
		assert.Len(t, player.rack, 2)
		assert.Len(t, game.sets, 1)
		assert.Len(t, game.sets[0].(*set).tiles, 3)
	})
	t.Run("ShouldReturnErrorOnBadCombine", func(t *testing.T) {
		run := &set{tiles: createRunTiles(t, 1, 4, ColorRed)}
		game := &instance{sets: []Set{run}}
		player := &player{rack: []Piece{NewPiece(Value(5), ColorRed), NewPiece(Value(7), ColorRed)}}
		readFromStringInput(t, "s0,3\nr0\nr1\ndone\n")
		err := player.combine(game)
		assert.Error(t, err)
		assert.Len(t, player.rack, 2)
		assert.Len(t, game.sets, 1)
		assert.Len(t, game.sets[0].(*set).tiles, 4)
	})
	t.Run("ShouldReturnErrorOnReaderError", func(t *testing.T) {
		player := new(player)
		game := new(instance)
		readerError(t)
		err := player.combine(game)
		assert.Error(t, err)
	})
}
