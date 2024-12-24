package model

import (
	"errors"
	"fmt"
)

var (
	ansiColors = map[Color]int{
		ColorBlack: 37,
		ColorBlue:  35,
		ColorRed:   31,
		ColorGreen: 32,
	}
)

type (
	Set interface {
		// Insert(p Piece) error
		Remove(index, position int, pn ...Piece) (Set, error)
		Split(piece Piece) (Set, error)
		Piece(index int) (Piece, error)
		String() string
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

func (s *set) String() string {
	var output string
	for i, t := range s.tiles {
		color := t.(*piece).Color
		if t.IsJoker() {
			output = output + fmt.Sprintf("[%d] (Joker)\n", i)
		} else {
			output = output + fmt.Sprintf("[%d] (\x1b[%dm%d\x1b[0m)\n", i, ansiColors[color], t.Value())
		}
	}
	return output
}

func (s *set) Piece(index int) (Piece, error) {
	if len(s.tiles) == 0 {
		return nil, errors.New(InvalidSet)
	}
	if index < 0 || index >= len(s.tiles) {
		return nil, errors.New(IndexOutOfBounds(len(s.tiles)-1))
	}
	return s.tiles[index], nil
}

func isGroup(s *set) bool {
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

func isRun(s *set) bool {
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

func (s *set) insertIntoGroup(piece Piece) error {
	if !piece.IsSameValue(s.tiles[0]) {
		return errors.New(WrongValueForGroup)
	}
	s.insertPiece(piece, len(s.tiles))
	return nil
}

func (s *set) insertIntoRun(piece Piece) error {
	if !piece.IsSameColor(s.tiles[0]) {
		return errors.New(WrongColorForRun)
	}
	if first := s.tiles[0]; piece.Value()+1 == first.Value() {
		s.insertPiece(piece, 0)
		return nil
	}
	if last := s.tiles[len(s.tiles)-1]; piece.Value()-1 == last.Value() {
		s.insertPiece(piece, len(s.tiles))
		return nil
	}
	return errors.New(CannotInsertIntoRun)
}

func (s *set) Insert(piece Piece) error {
	if !isValidPiece(piece) {
		return errors.New(InvalidPiece)
	}
	if isGroup(s) {
		return s.insertIntoGroup(piece)
	} else if isRun(s) {
		return s.insertIntoRun(piece)
	}
	return errors.New(InvalidSet)
}

func isValidSet(s *set) bool {
	if s == nil || len(s.tiles) < 3 {
		return false
	}
	return isGroup(s) || isRun(s)
}

func (s *set) removePiece(index int) {
	s.tiles = append(s.tiles[:index], s.tiles[index+1:]...)
}

func (s *set) cloneTiles() []Piece {
	clone := make([]Piece, len(s.tiles))
	copy(clone, s.tiles)
	return clone
}

func (s *set) Remove(index, position int, pn ...Piece) (Set, error) {
	if len(s.tiles) < 4 || len(pn)+1 < 3 {
		return nil, errors.New(TooFewPieces)
	}
	if index < 0 || index > len(s.tiles)-1 {
		return nil, errors.New(IndexOutOfBounds(len(s.tiles)-1))
	}
	if position < 0 || position > len(pn) {
		return nil, errors.New(IndexOutOfBounds(len(pn), "position"))
	}
	piece, newSet, original := s.tiles[index], &set{pn}, s.cloneTiles()
	newSet.insertPiece(piece, position)
	s.removePiece(index)
	if !isValidSet(newSet) || !isValidSet(s) {
		s.tiles = original
		return nil, errors.New(InvalidSet)
	}
	return newSet, nil
}

func (s *set) Split(piece Piece) (Set, error) {
	if !isValidSet(s) || !isRun(s) {
		return nil, errors.New(InvalidSet)
	}
	if len(s.tiles)+1 < 6 {
		return nil, errors.New(TooFewPieces)
	}
	for index, p := range s.tiles {
		minPieces := index > 1 && len(s.tiles)-index > 2
		if minPieces && piece.IsSameValue(p) {
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
