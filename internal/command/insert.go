package command

import (
	"errors"
	"lets-play-rummikub/internal/constants"
	"lets-play-rummikub/internal/model"
	"strings"
)

type insert struct {
	player     model.Player
	game       model.Game
	set        model.Set
	piece      model.Piece
	index      int
	undoGame   model.Game
	undoPlayer model.Player
}

func Insert(player model.Player, game model.Game, input string) (Command, error) {
	selections := strings.Split(input, " ")
	if len(selections) != 3 {
		return nil, errors.New(constants.TooFewArguments)
	}
	setSelection, pieceSelection, position := selections[0], selections[1], selections[2]
	index, err := parseInt(position)
	if err != nil {
		return nil, err
	}
	setIndex, err := parseInt(setSelection)
	if err != nil {
		return nil, err
	}
	pieceIndex, err := parseInt(pieceSelection[1:])
	if err != nil {
		return nil, err
	}
	selectPiece, err := getPieceFrom(pieceSelection[0], player, game)
	if err != nil {
		return nil, err
	}
	set, err := game.Set(setIndex)
	if err != nil {
		return nil, errors.New(constants.InvalidSetSelection)
	}
	piece, err := selectPiece.Piece(pieceIndex)
	if err != nil {
		return nil, errors.New(constants.InvalidPieceSelection)
	}
	return &insert{player, game, set, piece, index, nil, nil}, nil
}

func (i *insert) Undo() {
	i.game.Restore(i.undoGame)
	i.player.Restore(i.undoPlayer)
}

func (i *insert) Invoke() {
	if insert, err := i.set.Insert(i.piece, i.index); err != nil {
		i.game.Notify(err.Error())
	} else {
		i.undoGame = i.game.Clone()
		i.undoPlayer = i.player.Clone()
		i.player.RemovePiece(i.piece)
		i.game.RemovePieces(i.piece)
		i.game.ReplaceSet(i.set, insert)
	}
}
