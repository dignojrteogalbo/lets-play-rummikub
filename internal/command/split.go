package command

import (
	"errors"
	"lets-play-rummikub/internal/constants"
	"lets-play-rummikub/internal/model"
	"strings"
)

type split struct {
	game     model.Game
	set      model.Set
	index    int
	undoGame model.Game
}

func Split(game model.Game, input string) (Command, error) {
	selections := strings.Split(input, " ")
	if len(selections) != 2 {
		return nil, errors.New(constants.TooFewArguments)
	}
	setSelection, splitIndex := selections[0], selections[1]
	setIndex, err := parseInt(setSelection)
	if err != nil {
		return nil, err
	}
	set, err := game.Set(setIndex)
	if err != nil {
		return nil, err
	}
	index, err := parseInt(splitIndex)
	if err != nil {
		return nil, err
	}
	return &split{game, set, index, nil}, nil
}

func (s *split) Undo() {
	s.game.Restore(s.undoGame)
	s.game.Notify()
}

func (s *split) Invoke() {
	if lowerSet, upperSet, err := s.set.Split(s.index); err != nil {
		s.game.Notify(err.Error())
	} else {
		s.undoGame = s.game.Clone()
		s.game.ReplaceSet(s.set, lowerSet)
		s.game.AddSet(upperSet)
	}
}
