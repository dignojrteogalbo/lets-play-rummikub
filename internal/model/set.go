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
		Insert(p Piece, index int) (Set, error)
		Remove(index int) (Set, Piece, error)
		Split(index int) (Set, Set, error)
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
	if len(s.tiles) < 3 {
		return false
	}
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

func (s *set) findIndex(piece Piece) int {
	if s == nil || len(s.tiles) == 0 || !isValidPiece(piece) {
		return -1
	}
	for index, p := range s.tiles {
		if piece.IsSamePiece(p) {
			return index
		}
	}
	return -1
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

func (s *set) Insert(piece Piece, index int) (Set, error) {
	if len(s.tiles) != 0 && s.findIndex(piece) >= 0 {
		return nil, errors.New(InvalidPiece)
	}
	clone := &set{tiles: s.cloneTiles()}
	clone.insertPiece(piece, index)
	return clone, nil
}

//region remove set

func (s *set) Remove(index int) (Set, Piece, error) {
	if index < 0 || index >= len(s.tiles) {
		return nil, nil, errors.New(IndexOutOfBounds(len(s.tiles)-1))
	}
	if len(s.tiles) == 0 {
		return nil, nil, errors.New(InvalidSet)
	}
	clone := &set{tiles: s.cloneTiles()}
	piece := clone.tiles[index]
	clone.removePiece(index)
	return clone, piece, nil
}

//region split set

func (s *set) Split(index int) (Set, Set, error) {
	if index < 0 || index >= len(s.tiles) {
		return nil, nil, errors.New(IndexOutOfBounds(len(s.tiles)-1))
	}
	if len(s.tiles) < 2 {
		return nil, nil, errors.New(TooFewPieces)
	}
	clone := s.cloneTiles()
	return &set{tiles: clone[:index]}, &set{tiles: clone[index+1:]}, nil
}

//region combine set

func Combine(pieces ...Piece) Set {
	return &set{tiles: pieces}
}
