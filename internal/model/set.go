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
		Insert(p Piece) error
		Remove(p Piece) error
		Split(p Piece) (Set, error)
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

//region set validation

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

func isValidSet(s *set) bool {
	if s == nil || len(s.tiles) < 3 {
		return false
	}
	return isGroup(s) || isRun(s)
}

//region set manipulation

func (s *set) insertPiece(piece Piece, index int) {
	if index == len(s.tiles) {
		s.tiles = append(s.tiles, piece)
	} else {
		s.tiles = append(s.tiles[:index+1], s.tiles[index:]...)
		s.tiles[index] = piece
	}
}

func (s *set) removePiece(index int) {
	s.tiles = append(s.tiles[:index], s.tiles[index+1:]...)
}

func (s *set) cloneTiles() []Piece {
	clone := make([]Piece, len(s.tiles))
	copy(clone, s.tiles)
	return clone
}

//region insert set

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

//region remove set

func (s *set) Remove(piece Piece) error {
	if len(s.tiles) <= 3 {
		return errors.New(TooFewPieces)
	}
	clone := s.cloneTiles()
	index := -1
	for i, p := range s.tiles {
		if piece.IsSamePiece(p) {
			index = i
			break
		}
	}
	if index < 0 {
		return errors.New(InvalidPiece)
	}
	remove := &set{tiles: clone}
	remove.removePiece(index)
	if !isValidSet(remove) {
		return errors.New(InvalidSet)
	}
	s.tiles = remove.tiles
	return nil
}

//region split set

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

//region combine set

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
