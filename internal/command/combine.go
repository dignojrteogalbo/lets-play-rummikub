package command

import "lets-play-rummikub/internal/model"

type combine struct {
	player     model.Player
	game       model.Game
	pieces     []model.Piece
	undoGame   model.Game
	undoPlayer model.Player
}

func Combine(player model.Player, game model.Game, input string) (Command, error) {
	if pieces, err := parseSelectedPieces(input, player, game); err != nil {
		return nil, err
	} else {
		return &combine{player, game, pieces, nil, nil}, nil
	}
}

func (c *combine) Undo() {
	c.game.Restore(c.undoGame)
	c.player.Restore(c.undoPlayer)
	c.game.Notify()
}

func (c *combine) Invoke() {
	c.undoGame, c.undoPlayer = c.game.Clone(), c.player.Clone()
	set := model.Combine(c.pieces...)
	c.player.RemovePiece(c.pieces...)
	c.game.RemovePieces(c.pieces...)
	c.game.AddSet(set)
	c.game.Notify()
}
