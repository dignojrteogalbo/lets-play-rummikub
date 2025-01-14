package command

import (
	"lets-play-rummikub/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetName(t *testing.T) {
	game := model.NewGame(1)
	player := game.CurrentPlayer()
	command := SetName(player, "new name")
	assert.NotNil(t, command)
	result := command.(*setName)
	assert.Same(t, result.player, player)
	assert.Equal(t, result.name, "new name")
}

func TestInvokeSetName(t *testing.T) {
	game := model.NewGame(1)
	player := game.CurrentPlayer()
	command := SetName(player, "new name")
	assert.NotNil(t, command)
	command.Invoke()
	assert.Equal(t, player.Name(), "new name")
}

func TestUndoSetName(t *testing.T) {
	game := model.NewGame(1)
	player := game.CurrentPlayer()
	command := SetName(player, "new name")
	assert.NotNil(t, command)
	command.Invoke()
	command.Undo()
	assert.Equal(t, player.Name(), "new name")
}
