package command

import (
	"errors"
	"lets-play-rummikub/internal/constants"
	"lets-play-rummikub/internal/history"
	"lets-play-rummikub/internal/model"
	"strconv"
	"strings"
)

type Command interface {
	Invoke()
	history.Undoable
}

func parseInt(input string) (int, error) {
	result, err := strconv.ParseInt(input, 0, 16)
	if err != nil {
		return -1, errors.New(constants.InvalidNumberInput)
	}
	return int(result), nil
}

func getPieceFrom(from byte, options ...model.HasPiece) (model.HasPiece, error) {
	for _, opt := range options {
		switch opt.(type) {
		case model.Player:
			if from == 'r' {
				return opt, nil
			}
		case model.Game:
			if from == 'p' {
				return opt, nil
			}
		case model.Set:
			if from == 's' {
				return opt, nil
			}
		default:
			continue
		}
	}
	return nil, errors.New(constants.InvalidPieceSelection)
}

func parseSelectedPieces(input string, options ...model.HasPiece) ([]model.Piece, error) {
	pieces := make([]model.Piece, 0)
	for _, selection := range strings.Split(input, " ") {
		from, pieceIndex := selection[0], selection[1:]
		selectPiece, err := getPieceFrom(from, options...)
		if err != nil {
			return nil, err
		}
		index, err := parseInt(pieceIndex)
		if err != nil {
			return nil, err
		}
		if piece, err := selectPiece.Piece(index); err != nil {
			return nil, err
		} else {
			pieces = append(pieces, piece)
		}
	}
	return pieces, nil
}
