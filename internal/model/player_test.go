package model

import (
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

// func TestPlayerInsert(t *testing.T) {
// 	t.Run("ShouldInsertIntoSet", func(t *testing.T) {
// 		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
// 		game := &instance{board: []Set{existingSet}}
// 		player := &player{rack: []Piece{NewPiece(Value(4), ColorGreen)}}
// 		readFromStringInput(t, player, "0\nr0\n3\n")
// 		err := player.insert(game)
// 		assert.NoError(t, err)
// 		assert.Empty(t, player.rack)
// 		assert.Len(t, game.board, 1)
// 		assert.Len(t, game.board[0].(*set).tiles, existingSet.Len()+1)
// 		assert.NotSame(t, existingSet, game.board[0])
// 	})
// 	t.Run("ShouldReturnErrorOnBadSetSelection", func(t *testing.T) {
// 		set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
// 		game := &instance{board: []Set{set}}
// 		player := &player{rack: []Piece{NewPiece(Value(5), ColorGreen)}}
// 		readFromStringInput(t, player, "2\n")
// 		err := player.insert(game)
// 		assert.Error(t, err)
// 		assert.Len(t, player.rack, 1)
// 		assert.Len(t, set.tiles, 3)
// 	})
// 	t.Run("ShouldReturnErrorOnErrorPieceSelection", func(t *testing.T) {
// 		set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
// 		game := &instance{board: []Set{set}}
// 		player := &player{rack: []Piece{NewPiece(Value(5), ColorGreen)}}
// 		readFromStringInput(t, player, "0\n-1\n")
// 		err := player.insert(game)
// 		assert.Error(t, err)
// 		assert.Len(t, player.rack, 1)
// 		assert.Len(t, set.tiles, 3)
// 	})
// 	t.Run("ShouldReturnErrorOnBadPieceSelection", func(t *testing.T) {
// 		set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
// 		game := &instance{board: []Set{set}}
// 		player := &player{rack: []Piece{NewPiece(Value(5), ColorGreen)}}
// 		readFromStringInput(t, player, "0\n0\n")
// 		err := player.insert(game)
// 		assert.Error(t, err)
// 		assert.Len(t, player.rack, 1)
// 		assert.Len(t, set.tiles, 3)
// 	})
// 	t.Run("ShouldReturnErrorOnErrorIndexSelection", func(t *testing.T) {
// 		set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
// 		game := &instance{board: []Set{set}}
// 		player := &player{rack: []Piece{NewPiece(Value(5), ColorGreen)}}
// 		readFromStringInput(t, player, "0\nr0\n")
// 		err := player.insert(game)
// 		assert.Error(t, err)
// 		assert.Len(t, player.rack, 1)
// 		assert.Len(t, set.tiles, 3)
// 	})
// 	t.Run("ShouldReturnErrorOnBadIndexSelection", func(t *testing.T) {
// 		set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
// 		game := &instance{board: []Set{set}}
// 		player := &player{rack: []Piece{NewPiece(Value(5), ColorGreen)}}
// 		readFromStringInput(t, player, "0\nr0\nbad\n")
// 		err := player.insert(game)
// 		assert.Error(t, err)
// 		assert.Len(t, player.rack, 1)
// 		assert.Len(t, set.tiles, 3)
// 	})
// 	t.Run("ShouldReturnErrorOnBadInsert", func(t *testing.T) {
// 		set := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
// 		game := &instance{board: []Set{set}}
// 		player := &player{rack: []Piece{NewPiece(Value(5), ColorGreen)}}
// 		readFromStringInput(t, player, "0\nr0\n-1\n")
// 		err := player.insert(game)
// 		assert.Error(t, err)
// 		assert.Len(t, player.rack, 1)
// 		assert.Len(t, set.tiles, 3)
// 	})
// }

// func TestPlayerCombine(t *testing.T) {
// 	t.Run("ShouldCombinePiecesFromRack", func(t *testing.T) {
// 		game := new(instance)
// 		player := &player{rack: createRunTiles(t, 1, 5, ColorBlue)}
// 		readFromStringInput(t, player, "r0\nr1\nr2\nr3\nr4\ndone\n")
// 		err := player.combine(game)
// 		assert.NoError(t, err)
// 		assert.Empty(t, player.rack)
// 		assert.Len(t, game.board, 1)
// 		assert.Len(t, game.board[0].(*set).tiles, 5)
// 	})
// 	t.Run("ShouldCombinePiecesFromRackAndBoard", func(t *testing.T) {
// 		game := &instance{loose: createRunTiles(t, 1, 4, ColorRed)}
// 		player := &player{rack: []Piece{NewPiece(Value(5), ColorRed), NewPiece(Value(6), ColorRed)}}
// 		readFromStringInput(t, player, "p3\nr0\nr1\ndone\n")
// 		err := player.combine(game)
// 		assert.NoError(t, err)
// 		assert.Empty(t, player.rack)
// 		assert.Len(t, game.board, 1)
// 		assert.Len(t, game.board[0].(*set).tiles, 3)
// 		assert.Len(t, game.loose, 3)
// 	})
// 	t.Run("ShouldReturnErrorOnBadPiece", func(t *testing.T) {
// 		game := &instance{loose: createRunTiles(t, 1, 3, ColorRed)}
// 		player := &player{rack: []Piece{NewPiece(Value(4), ColorRed), NewPiece(Value(5), ColorRed)}}
// 		readFromStringInput(t, player, "r0\nr1\np3\ndone\n")
// 		err := player.combine(game)
// 		assert.Error(t, err)
// 		assert.Len(t, player.rack, 2)
// 		assert.Len(t, game.loose, 3)
// 	})
// }

// func TestPlayerRemove(t *testing.T) {
// 	t.Run("ShouldRemoveFromSet", func(t *testing.T) {
// 		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
// 		game := &instance{board: []Set{existingSet}}
// 		player := new(player)
// 		readFromStringInput(t, player, "0\n0\n")
// 		err := player.remove(game)
// 		assert.NoError(t, err)
// 		assert.Len(t, game.board, 1)
// 		assert.NotSame(t, game.board[0], existingSet)
// 		assert.Len(t, game.board[0].(*set).tiles, existingSet.Len()-1)
// 		assert.Len(t, game.loose, 1)
// 		assert.Same(t, existingSet.tiles[0], game.loose[0])
// 	})
// 	t.Run("ShouldReturnErrorOnBadSetSelection", func(t *testing.T) {
// 		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
// 		game := &instance{board: []Set{existingSet}}
// 		player := new(player)
// 		readFromStringInput(t, player, "1\n")
// 		err := player.remove(game)
// 		assert.Error(t, err)
// 		assert.Len(t, game.board, 1)
// 		assert.Same(t, game.board[0], existingSet)
// 		assert.Len(t, existingSet.tiles, 3)
// 	})
// 	t.Run("ShouldReturnErrorOnReaderError", func(t *testing.T) {
// 		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
// 		game := &instance{board: []Set{existingSet}}
// 		player := new(player)
// 		readFromStringInput(t, player, "0\n")
// 		err := player.remove(game)
// 		assert.Error(t, err)
// 		assert.Len(t, game.board, 1)
// 		assert.Same(t, game.board[0], existingSet)
// 		assert.Len(t, existingSet.tiles, 3)
// 	})
// 	t.Run("ShouldReturnErrorOnBadNumberInput", func(t *testing.T) {
// 		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
// 		game := &instance{board: []Set{existingSet}}
// 		player := new(player)
// 		readFromStringInput(t, player, "0\nbad input\n")
// 		err := player.remove(game)
// 		assert.Error(t, err)
// 		assert.Len(t, game.board, 1)
// 		assert.Same(t, game.board[0], existingSet)
// 		assert.Len(t, existingSet.tiles, 3)
// 	})
// 	t.Run("ShouldReturnErrorOnBadPieceSelection", func(t *testing.T) {
// 		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
// 		game := &instance{board: []Set{existingSet}}
// 		player := new(player)
// 		readFromStringInput(t, player, "0\n-1\n")
// 		err := player.remove(game)
// 		assert.Error(t, err)
// 		assert.Len(t, game.board, 1)
// 		assert.Same(t, game.board[0], existingSet)
// 		assert.Len(t, existingSet.tiles, 3)
// 	})
// 	t.Run("ShouldReturnErrorOnBadRemove", func(t *testing.T) {
// 		existingSet := &set{tiles: []Piece{nil}}
// 		game := &instance{board: []Set{existingSet}}
// 		player := new(player)
// 		readFromStringInput(t, player, "0\n0\n")
// 		err := player.remove(game)
// 		assert.Error(t, err)
// 		assert.Len(t, game.board, 1)
// 		assert.Same(t, game.board[0], existingSet)
// 		assert.Len(t, existingSet.tiles, 1)
// 	})
// }

// func TestPlayerSplit(t *testing.T) {
// 	t.Run("ShouldSplitSet", func(t *testing.T) {
// 		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
// 		game := &instance{board: []Set{existingSet}}
// 		player := new(player)
// 		readFromStringInput(t, player, "0\n1\n")
// 		err := player.split(game)
// 		assert.NoError(t, err)
// 		assert.Len(t, game.board, 2)
// 		assert.NotSame(t, game.board[0], existingSet)
// 		assert.NotSame(t, game.board[1], existingSet)
// 		assert.Len(t, game.board[0].(*set).tiles, 1)
// 		assert.Len(t, game.board[1].(*set).tiles, 2)
// 	})
// 	t.Run("ShouldReturnErrorOnBadSetSelection", func(t *testing.T) {
// 		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
// 		game := &instance{board: []Set{existingSet}}
// 		player := new(player)
// 		readFromStringInput(t, player, "1\n")
// 		err := player.split(game)
// 		assert.Error(t, err)
// 		assert.Len(t, game.board, 1)
// 		assert.Same(t, game.board[0], existingSet)
// 		assert.Len(t, existingSet.tiles, 3)
// 	})
// 	t.Run("ShouldReturnErrorOnCannotSplitSet", func(t *testing.T) {
// 		existingSet := &set{tiles: []Piece{NewPiece(Value(7), ColorBlack)}}
// 		game := &instance{board: []Set{existingSet}}
// 		player := new(player)
// 		readFromStringInput(t, player, "0\n")
// 		err := player.split(game)
// 		assert.Error(t, err)
// 		assert.Len(t, game.board, 1)
// 		assert.Same(t, game.board[0], existingSet)
// 		assert.Len(t, existingSet.tiles, 1)
// 	})
// 	t.Run("ShouldReturnErrorOnReaderError", func(t *testing.T) {
// 		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
// 		game := &instance{board: []Set{existingSet}}
// 		player := new(player)
// 		readFromStringInput(t, player, "0\n")
// 		err := player.split(game)
// 		assert.Error(t, err)
// 		assert.Len(t, game.board, 1)
// 		assert.Same(t, game.board[0], existingSet)
// 		assert.Len(t, existingSet.tiles, 3)
// 	})
// 	t.Run("ShouldReturnErrorOnBadIndexInput", func(t *testing.T) {
// 		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
// 		game := &instance{board: []Set{existingSet}}
// 		player := new(player)
// 		readFromStringInput(t, player, "0\nbad input\n")
// 		err := player.split(game)
// 		assert.Error(t, err)
// 		assert.Len(t, game.board, 1)
// 		assert.Same(t, game.board[0], existingSet)
// 		assert.Len(t, existingSet.tiles, 3)
// 	})
// 	t.Run("ShouldReturnErrorOnBadSplit", func(t *testing.T) {
// 		existingSet := &set{tiles: createRunTiles(t, 1, 3, ColorGreen)}
// 		game := &instance{board: []Set{existingSet}}
// 		player := new(player)
// 		readFromStringInput(t, player, "0\n3\n")
// 		err := player.split(game)
// 		assert.Error(t, err)
// 		assert.Len(t, game.board, 1)
// 		assert.Same(t, game.board[0], existingSet)
// 		assert.Len(t, existingSet.tiles, 3)
// 	})
// }
