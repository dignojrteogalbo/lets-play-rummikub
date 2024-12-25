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
		String() string
		Piece(index int) (Piece, error)
		Clone() Set
	}

	set struct {
		tiles []Piece
	}
)

func (s *set) String() string {
	var output string
	for i, t := range s.tiles {
		if t.IsJoker() {
			output = output + fmt.Sprintf("[%d] (Joker)\n", i)
		} else {
			output = output + fmt.Sprintf("[%d] (\x1b[%dm%d\x1b[0m)\n", i, ansiColors[t.Color()], t.Value())
		}
	}
	return output
}

func (s *set) Piece(index int) (Piece, error) {
	if len(s.tiles) == 0 {
		return nil, errors.New(InvalidSet)
	}
	if index < 0 || index >= len(s.tiles) {
		return nil, errors.New(IndexOutOfBounds(len(s.tiles) - 1))
	}
	return s.tiles[index], nil
}

func (s *set) Clone() Set {
	return &set{tiles: s.cloneTiles()}
}

//region set validation

func isGroup(s *set) bool {
	if len(s.tiles) < 3 || len(s.tiles) > 4 {
		return false
	}
	startPiece := s.tiles[0]
	colors := map[Color]bool{
		ColorBlack: false,
		ColorBlue:  false,
		ColorRed:   false,
		ColorGreen: false,
	}
	totalColors := 0
	for _, piece := range s.tiles {
		if piece.IsJoker() {
			totalColors++
			continue
		}
		if !startPiece.IsSameValue(piece) {
			return false
		}
		if !colors[piece.Color()] {
			colors[piece.Color()] = true
			totalColors++
		}
	}
	return totalColors > 2
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
	if len(s.tiles) == 4 {
		return errors.New(CannotInsert)
	}
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
	return errors.New(CannotInsert)
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

func (s *set) splitSet(piece Piece, index int) Set {
	s.insertPiece(piece, index)
	clone := s.cloneTiles()
	s.tiles = clone[:index+1]
	return &set{tiles: clone[index+1:]}
}

func (s *set) Split(piece Piece) (Set, error) {
	if !isRun(s) {
		return nil, errors.New(CannotSplit)
	}
	if len(s.tiles) < 5 {
		return nil, errors.New(TooFewPieces)
	}
	index := -1
	for i, p := range s.tiles {
		if piece.IsSameValue(p) && piece.IsSameColor(p) {
			index = i
			break
		}
	}
	if index < 0 {
		return nil, errors.New(InvalidPiece)
	} else if index < 2 || len(s.tiles)-index+1 < 3 {
		return nil, errors.New(TooFewPieces)
	}
	return s.splitSet(piece, index), nil
}

//region combine set

func Combine(pieces ...Piece) (Set, error) {
	if len(pieces) < 3 {
		return nil, errors.New(TooFewPieces)
	}
	combined := &set{tiles: pieces}
	if !isValidSet(combined) {
		return nil, errors.New(InvalidSet)
	}
	return combined, nil
}
