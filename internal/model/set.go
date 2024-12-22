package model

import (
	"errors"
)

type (
	Set interface {
		Insert(p Piece, index int) error
		Remove(index, position int, pn ...Piece) (Set, error)
		IsGroup() bool
		IsRun() bool
	}

	set struct {
		tiles []Piece
	}

	combineArguments struct {
		fromSet   *set
		atIndex   int
		takePiece *piece
	}
)

func (s *set) IsGroup() bool {
	startPiece := s.tiles[0]
	for _, piece := range s.tiles {
		if piece.IsJoker() {
			continue
		}
		if !startPiece.IsSameValue(piece) {
			return false
		}
	}
	return true
}

func (s *set) IsRun() bool {
	startPiece := s.tiles[0]
	var previous, current Piece
	for _, piece := range s.tiles {
		previous, current = current, piece
		if current.IsJoker() || previous == nil || previous.IsJoker() {
			continue
		}
		if !startPiece.IsSameColor(piece) {
			return false
		}
		if previous.Value()+1 != current.Value() {
			return false
		}
	}
	return true
}

func (s *set) insertPiece(piece Piece, index int) {
	if index == len(s.tiles) {
		s.tiles = append(s.tiles, piece)
	} else {
		s.tiles = append(s.tiles[:index+1], s.tiles[index:]...)
		s.tiles[index] = piece
	}
}

func (s *set) insertIntoGroup(piece Piece, index int) error {
	if !piece.IsSameValue(s.tiles[0]) {
		return errors.New(WrongValueForGroup)
	}
	s.insertPiece(piece, index)
	return nil
}

func (s *set) insertIntoRun(piece Piece, index int) error {
	if !piece.IsSameColor(s.tiles[0]) {
		return errors.New(WrongColorForRun)
	}
	first, last := s.tiles[0], s.tiles[len(s.tiles)-1]
	validStartPiece := piece.Value()+1 == first.Value()
	validEndPiece := piece.Value()-1 == last.Value()
	if (index == 0 && validStartPiece) || (index == len(s.tiles) && validEndPiece) {
		s.insertPiece(piece, index)
		return nil
	}
	return errors.New(CannotInsertIntoRun)
}

func (s *set) Insert(piece Piece, index int) error {
	if !isValidPiece(piece) {
		return errors.New(InvalidPiece)
	}
	if index < 0 || index > len(s.tiles) {
		return errors.New(IndexOutOfBounds(len(s.tiles)))
	}
	if s.IsGroup() {
		return s.insertIntoGroup(piece, index)
	}
	if s.IsRun() {
		return s.insertIntoRun(piece, index)
	}
	return errors.New(InvalidSet)
}

func isValidSet(s Set) bool {
	if s == nil {
		return false
	}
	if set, ok := s.(*set); !ok || len(set.tiles) < 3 {
		return false
	}
	return s.IsGroup() || s.IsRun()
}

func (s *set) removePiece(index int) {
	s.tiles = append(s.tiles[:index], s.tiles[index+1:]...)
}

func (s *set) Remove(index, position int, pn ...Piece) (Set, error) {
	if len(s.tiles) < 4 || len(pn)+1 < 3 {
		return nil, errors.New(TooFewPieces)
	}
	if index < 0 || index >= len(s.tiles) {
		return nil, errors.New(IndexOutOfBounds(len(s.tiles)))
	}
	if position < 0 || position >= len(pn) {
		return nil, errors.New(IndexOutOfBounds(len(pn), "position"))
	}
	piece, newSet := s.tiles[index], &set{pn}
	newSet.insertPiece(piece, position)
	if !isValidSet(newSet) {
		return nil, errors.New(InvalidSet)
	}
	s.removePiece(index)
	return newSet, nil
}

func (s *set) Split(piece Piece) (Set, error) {
	if !isValidSet(s) || !s.IsRun() {
		return nil, errors.New(InvalidSet)
	}
	if len(s.tiles)+1 < 6 {
		return nil, errors.New(TooFewPieces)
	}
	for index, p := range s.tiles {
		if index > 1 && piece.IsSameValue(p) {
			s.insertPiece(piece, index)
			splitSet := &set{s.tiles[index+1:]}
			s.tiles = s.tiles[:index+1]
			return splitSet, nil
		}
	}
	return nil, errors.New(CannotSplitRun)
}

func processCombineArguments(setsAndIndexesOrPieces ...any) ([]combineArguments, error) {
	combination := make([]combineArguments, 0)
	for len(setsAndIndexesOrPieces) > 0 {
		if takePiece, ok := setsAndIndexesOrPieces[0].(*piece); ok {
			combination = append(combination, combineArguments{nil, 0, takePiece})
			setsAndIndexesOrPieces = setsAndIndexesOrPieces[1:]
			continue
		}
		if len(setsAndIndexesOrPieces) < 2 {
			return nil, errors.New(InvalidCombineArguments)
		}
		takeSet, ok := setsAndIndexesOrPieces[0].(*set)
		if !ok {
			return nil, errors.New(InvalidCombineArguments)
		}
		takePiece, ok := setsAndIndexesOrPieces[1].(int)
		if !ok {
			return nil, errors.New(InvalidCombineArguments)
		}
		combination = append(combination, combineArguments{takeSet, takePiece, nil})
		setsAndIndexesOrPieces = setsAndIndexesOrPieces[2:]
	}
	if len(combination) < 3 {
		return nil, errors.New(TooFewPieces)
	}
	return combination, nil
}

func isValidCombination(combination []combineArguments) (Set, bool) {
	tiles := make([]Piece, 0)
	for _, c := range combination {
		if c.takePiece != nil {
			tiles = append(tiles, c.takePiece)
			continue
		}
		tiles = append(tiles, c.fromSet.tiles[c.atIndex])
	}
	newSet := &set{tiles}
	if isValidSet(newSet) {
		return newSet, true
	}
	return nil, false
}

func applyCombination(combination []combineArguments) error {
	revert := make([][]Piece, 0)
	for index, c := range combination {
		if c.takePiece != nil {
			continue
		}
		original := make([]Piece, len(c.fromSet.tiles))
		copy(original, c.fromSet.tiles)
		revert = append(revert, original)
		c.fromSet.removePiece(c.atIndex)
		if !isValidSet(c.fromSet) {
			for len(revert) > 0 {
				if combination[index].takePiece != nil {
					index--
					continue
				}
				combination[index].fromSet.tiles = revert[len(revert)-1]
				revert = revert[:len(revert)-1]
				index--
			}
			return errors.New(InvalidSet)
		}
	}
	return nil
}

func Combine(setsAndIndexesOrPieces ...any) (Set, error) {
	combination, err := processCombineArguments(setsAndIndexesOrPieces...)
	if err != nil {
		return nil, err
	}
	combineSet, valid := isValidCombination(combination)
	if !valid {
		return nil, errors.New(InvalidSet)
	}
	if err := applyCombination(combination); err != nil {
		return nil, err
	}
	return combineSet, nil
}
