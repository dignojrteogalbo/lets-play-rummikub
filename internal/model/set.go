package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"lets-play-rummikub/internal/constants"
)

type (
	Set interface {
		Len() int
		Size() int
		NumberOfJokers() int
		Insert(p Piece, index int) (Set, error)
		Remove(p Piece) (Set, error)
		Split(index int) (Set, Set, error)
		String() string
		MarshalJSON() ([]byte, error)
		Clone() Set
		IsValidSet() bool
		HasPiece
	}

	set struct {
		tiles []Piece
	}
)

func (s *set) MarshalJSON() ([]byte, error) {
	output := struct {
		Pieces []Piece `json:"pieces"`
	}{
		s.tiles,
	}
	return json.Marshal(output)
}

func (s *set) Len() int {
	if s == nil {
		return 0
	}
	return len(s.tiles)
}

func (s *set) Size() int {
	sum := Value(0)
	for _, piece := range s.tiles {
		sum = sum + piece.Value()
	}
	return int(sum)
}

func (s *set) NumberOfJokers() int {
	numberOfJokers := 0
	for _, piece := range s.tiles {
		if piece.IsJoker() {
			numberOfJokers = numberOfJokers + 1
		}
	}
	return numberOfJokers
}

func (s *set) String() string {
	var output string
	for i, t := range s.tiles {
		output = output + fmt.Sprintf("[%d] %s\n", i, t.String())
	}
	return output
}

func (s *set) Piece(index int) (Piece, error) {
	if index < 0 || index >= len(s.tiles) {
		return nil, errors.New(constants.InvalidPieceSelection)
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
	var startPiece Piece
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
		} else if startPiece == nil {
			startPiece = piece
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
	expectedColor, startingValue := Color(0), Value(0)
	for index, piece := range s.tiles {
		if piece.IsJoker() {
			continue
		}
		if expectedColor == Color(0) {
			expectedColor = piece.Color()
		} else if piece.Color() != expectedColor {
			return false
		}
		if startingValue == Value(0) {
			startingValue = piece.Value() - Value(index)
		} else if piece.Value() != startingValue+Value(index) {
			return false
		}
	}
	return true
}

func (s *set) IsValidSet() bool {
	if s == nil || len(s.tiles) < 3 || len(s.tiles) > 13 {
		return false
	}
	if s.NumberOfJokers() > 1 {
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
	if index < 0 || index > len(s.tiles) {
		return nil, errors.New(constants.IndexOutOfBounds(-1, len(s.tiles)+1))
	}
	if len(s.tiles) != 0 && s.findIndex(piece) >= 0 {
		return nil, errors.New(constants.InvalidPiece)
	}
	clone := &set{tiles: s.cloneTiles()}
	clone.insertPiece(piece, index)
	return clone, nil
}

//region remove set

func (s *set) Remove(piece Piece) (Set, error) {
	if len(s.tiles) == 0 {
		return nil, errors.New(constants.InvalidSet)
	}
	var index int
	if index = s.findIndex(piece); index < 0 {
		return nil, errors.New(constants.InvalidPiece)
	}
	clone := &set{tiles: s.cloneTiles()}
	clone.removePiece(index)
	return clone, nil
}

//region split set

func (s *set) Split(index int) (Set, Set, error) {
	if len(s.tiles) < 2 {
		return nil, nil, errors.New(constants.TooFewPieces)
	}
	if index < 1 || index >= len(s.tiles) {
		return nil, nil, errors.New(constants.IndexOutOfBounds(0, len(s.tiles)))
	}
	clone := s.cloneTiles()
	return &set{tiles: clone[:index]}, &set{tiles: clone[index:]}, nil
}

//region combine set

func Combine(pieces ...Piece) Set {
	return &set{tiles: pieces}
}
