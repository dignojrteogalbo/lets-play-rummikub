package model

import (
	"errors"
	"fmt"
)

type (
	Set interface {
		Insert(p Piece, index int) error
		Swap(p Piece, s Set, origin, destination int) error
		Split(index int) (Set, Set)
		IsGroup() bool
		IsRun() bool
	}

	set struct {
		tiles []Piece
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
		return errors.New("piece does not match the value of the group")
	}
	s.insertPiece(piece, index)
	return nil
}

func (s *set) insertIntoRun(piece Piece, index int) error {
	if !piece.IsSameColor(s.tiles[0]) {
		return errors.New("piece does not match the color of the run")
	}
	first, last := s.tiles[0], s.tiles[len(s.tiles)-1]
	validStartPiece := piece.Value()+1 == first.Value()
	validEndPiece := piece.Value()-1 == last.Value()
	if (index == 0 && validStartPiece) || (index == len(s.tiles) && validEndPiece) {
		s.insertPiece(piece, index)
		return nil
	}
	return fmt.Errorf("piece cannot be inserted into run")
}

func (s *set) Insert(piece Piece, index int) error {
	if !isValidPiece(piece) {
		return fmt.Errorf("piece is invalid")
	}
	if index < 0 || index > len(s.tiles) {
		return fmt.Errorf("index must be >= 0 and <= %d", len(s.tiles))
	}
	if s.IsGroup() {
		return s.insertIntoGroup(piece, index)
	}
	if s.IsRun() {
		return s.insertIntoRun(piece, index)
	}
	return errors.New("set is invalid")
}
