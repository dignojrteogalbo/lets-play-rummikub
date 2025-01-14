package command

import "lets-play-rummikub/internal/model"

type setName struct {
	player model.Player
	name   string
}

func SetName(player model.Player, name string) Command {
	return &setName{player, name}
}

func (*setName) Undo() {
	// not undoable
}

func (n *setName) Invoke() {
	n.player.SetName(n.name)
}
