package command

import (
	"errors"
	"lets-play-rummikub/internal/constants"
	"lets-play-rummikub/internal/model"
	"strings"
)

type remove struct {
	game     model.Game
	set      model.Set
	piece    model.Piece
	undoGame model.Game
}

func Remove(game model.Game, input string) (Command, error) {
	selections := strings.Split(input, " ")
	if len(selections) != 2 {
		return nil, errors.New(constants.TooFewArguments)
	}
	setSelection, pieceSelection := selections[0], selections[1]
	setIndex, err := parseInt(setSelection)
	if err != nil {
		return nil, err
	}
	set, err := game.Set(setIndex)
	if err != nil {
		return nil, err
	}
	pieceIndex, err := parseInt(pieceSelection)
	if err != nil {
		return nil, err
	}
	piece, err := set.Piece(pieceIndex)
	if err != nil {
		return nil, err
	}
	return &remove{game, set, piece, nil}, nil
}

func (r *remove) Undo() {
	r.game.Restore(r.undoGame)
	r.game.Notify()
}

func (r *remove) Invoke() {
	if remove, err := r.set.Remove(r.piece); err != nil {
		r.game.Notify(err.Error())
	} else {
		r.undoGame = r.game.Clone()
		r.game.AddLoosePiece(r.piece)
		r.game.ReplaceSet(r.set, remove)
		r.game.Notify()
	}
}
