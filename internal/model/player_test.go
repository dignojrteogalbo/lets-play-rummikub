package model

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	game := &instance{sets: []Set{&set{createGroupTiles(t, 4, Value(11))}}}
	player := &player{rack: []Piece{NewPiece(Value(7), ColorGreen), NewPiece(Value(3), ColorBlue)}}
	t.Run("ShouldSelectPieceFromRack", func(t *testing.T) {
		selected := player.selectPiece("r0\n", game)
		assert.NotNil(t, selected)
		assert.Same(t, selected.(*piece), player.rack[0].(*piece))
	})
	t.Run("ShouldReturnNilOnInvalidRackInput", func(t *testing.T) {
		selected := player.selectPiece("rbad input\n", game)
		assert.Nil(t, selected)
	})
	t.Run("ShouldReturnNilOnInvalidRackIndex", func(t *testing.T) {
		selected := player.selectPiece("r73\n", game)
		assert.Nil(t, selected)
	})
	t.Run("ShouldSelectPieceFromSet", func(t *testing.T) {
		selected := player.selectPiece("s0,2\n", game)
		assert.NotNil(t, selected)
		assert.Same(t, selected.(*piece), game.sets[0].(*set).tiles[2])
	})
	t.Run("ShouldReturnNilOnInvalidSelection", func(t *testing.T) {
		selected := player.selectPiece("s02\n", game)
		assert.Nil(t, selected)
	})
	t.Run("ShouldReturnNilOnBadSetInput", func(t *testing.T) {
		selected := player.selectPiece("sbad set,2\n", game)
		assert.Nil(t, selected)
	})
	t.Run("ShouldReturnNilOnBadPieceInput", func(t *testing.T) {
		selected := player.selectPiece("s0,bad piece\n", game)
		assert.Nil(t, selected)
	})
	t.Run("ShouldReturnNilOnInvalidSetSelection", func(t *testing.T) {
		selected := player.selectPiece("s2,2\n", game)
		assert.Nil(t, selected)
	})
	t.Run("ShouldReturnNilOnInvalidPieceSelection", func(t *testing.T) {
		selected := player.selectPiece("s0,7\n", game)
		assert.Nil(t, selected)
	})
	t.Run("ShouldReturnNilOnBadInput", func(t *testing.T) {
		selected := player.selectPiece("bad input\n", game)
		assert.Nil(t, selected)
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
