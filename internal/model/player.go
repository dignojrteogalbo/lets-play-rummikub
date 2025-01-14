package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"lets-play-rummikub/internal/constants"
)

type (
	Player interface {
		DealPiece(Piece)
		RackLen() int
		Score() uint16
		MarshalJSON() ([]byte, error)
		Name() string
		SetName(string)
		RemovePiece(pieces ...Piece)
		HasPiece
		Clone() Player
		Restore(Player)
	}

	player struct {
		name    string
		rack    []Piece
		endTurn chan bool
	}
)

func (p *player) MarshalJSON() ([]byte, error) {
	output := struct {
		Rack []Piece `json:"rack"`
	}{
		p.rack,
	}
	return json.Marshal(output)
}

func (p *player) Clone() Player {
	rack := make([]Piece, len(p.rack))
	copy(rack, p.rack)
	return &player{rack: rack}
}

func (p *player) Restore(restorePlayer Player) {
	restore, ok := restorePlayer.(*player)
	if !ok {
		return
	}
	p.rack = restore.rack
}

func NewPlayer() Player {
	return &player{"", make([]Piece, 0), make(chan bool)}
}

func (p *player) Name() string {
	if len(p.name) == 0 {
		return "Player"
	}
	return p.name
}

func (p *player) SetName(name string) {
	p.name = name
}

func (p *player) DealPiece(piece Piece) {
	if piece == nil {
		return
	}
	p.rack = append(p.rack, piece)
}

func (p *player) printRack() {
	rack := &set{tiles: p.rack}
	fmt.Printf("=== Rack ===\n%s=== Rack ===\n", rack.String())
}

func (p *player) Score() uint16 {
	score := uint16(0)
	for _, piece := range p.rack {
		score = score + uint16(piece.Value())
	}
	return score
}

func (player *player) RemovePiece(pieces ...Piece) {
	for _, piece := range pieces {
		for index, p := range player.rack {
			if piece.IsSamePiece(p) {
				player.rack = append(player.rack[:index], player.rack[index+1:]...)
				break
			}
		}
	}
}

func (p *player) Piece(index int) (Piece, error) {
	if index < 0 || index >= len(p.rack) {
		return nil, errors.New(constants.InvalidPieceSelection)
	}
	return p.rack[index], nil
}

func (p *player) RackLen() int {
	return len(p.rack)
}
