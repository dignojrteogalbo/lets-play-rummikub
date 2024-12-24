package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type (
	Player interface {
		DealPiece(Piece)
		StartTurn(game Game)
		Score() uint16
	}

	player struct {
		rack       []Piece
		totalMoves uint16
	}
)

func NewPlayer() Player {
	return &player{rack: make([]Piece, 0)}
}

func (p *player) DealPiece(piece Piece) {
	p.rack = append(p.rack, piece)
}

func (p *player) printRack() {
	rack := &set{tiles: p.rack}
	fmt.Printf("=== Rack ===\n%s=== Rack ===\n", rack.String())
}

func (player *player) StartTurn(game Game) {
	successfulMoves := uint16(0)
	for {
		player.printRack()
		input, err := Reader.ReadString('\n')
		command := strings.TrimSpace(input)
		if err != nil || input == "done" {
			break
		}
		parseCommand(command, &successfulMoves, player, game)
	}
}

func (p *player) Score() uint16 {
	score := uint16(0)
	for _, piece := range p.rack {
		score = score + uint16(piece.Value())
	}
	return score
}

func (p *player) TotalMoves() uint16 {
	return p.totalMoves
}

func (player *player) removePiece(piece Piece) {
	for index, p := range player.rack {
		if piece.IsSamePiece(p) {
			player.rack = append(player.rack[:index], player.rack[index+1:]...)
			return
		}
	}
}

func (p *player) promptForPiece(set Set) (Piece, error) {
	if set != nil {
		fmt.Println(set.String())
		fmt.Print("Select piece from set: ")
	} else {
		fmt.Print("Select piece from rack: ")
	}
	input, err := Reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	input = strings.TrimSpace(input)
	pieceIndex, err := parseInt(input)
	if err != nil {
		return nil, err
	}
	var piece Piece
	if set != nil {
		piece, err = set.Piece(pieceIndex)
	} else {
		piece, err = p.Piece(pieceIndex)
	}
	if err != nil {
		return nil, err
	}
	return piece, nil
}

func (p *player) promptForSet(game Game) (Set, error) {
	fmt.Print("Select set to insert: ")
	input, err := Reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	input = strings.TrimSpace(input)
	setIndex, err := parseInt(input)
	if err != nil {
		return nil, err
	}
	set, err := game.Set(setIndex)
	if err != nil {
		return nil, err
	}
	return set, nil
}

func (p *player) insert(game Game) error {
	set, err := p.promptForSet(game)
	if err != nil {
		return err
	}
	piece, err := p.promptForPiece(nil)
	if err != nil {
		return err
	}
	if err := set.Insert(piece); err != nil {
		return err
	}
	p.removePiece(piece)
	return nil
}

func parseCommand(input string, successfulMoves *uint16, player *player, game Game) {
	switch input {
	case "combine":

	case "insert":
		if err := player.insert(game); err != nil {
			fmt.Println(err)
		} else {
			*successfulMoves++
		}
	case "split":

	case "help":
		fmt.Println("combine <r#|s#,#> <r#|s#,#> ... <r#|s#,#>")
		fmt.Println("insert <set> <piece>")
		fmt.Println("split <piece> <set>")
	default:
		fmt.Println("invalid command")
	}
}

func parseInt(input string) (int, error) {
	result, err := strconv.ParseInt(input, 0, 16)
	if err != nil {
		return -1, errors.New(InvalidNumberInput)
	}
	return int(result), nil
}

func (p *player) Piece(index int) (Piece, error) {
	if index < 0 || index > len(p.rack) {
		return nil, errors.New(IndexOutOfBounds(len(p.rack)-1, "piece"))
	}
	return p.rack[index], nil
}

func (player *player) selectPiece(selection string, game Game) Piece {
	selection = strings.TrimSpace(selection)
	start := string(selection[0])
	if start == "r" {
		index, err := parseInt(selection[1:])
		if err != nil {
			fmt.Println("Invalid input")
			return nil
		}
		if piece, err := player.Piece(index); err != nil {
			fmt.Println("Invalid piece selection")
			return nil
		} else {
			return piece
		}
	}
	if start == "s" {
		tuple := strings.Split(selection[1:], ",")
		if len(tuple) != 2 {
			fmt.Println("Invalid input")
			return nil
		}
		setIndex, err := parseInt(tuple[0])
		if err != nil {
			fmt.Println("Invalid input")
			return nil
		}
		pieceIndex, err := parseInt(tuple[1])
		if err != nil {
			fmt.Println("Invalid input")
			return nil
		}
		set, err := game.Set(setIndex)
		if err != nil {
			fmt.Println("Invalid set selection")
			return nil
		}
		piece, err := set.Piece(pieceIndex)
		if err != nil {
			fmt.Println("Invalid piece selection")
			return nil
		}
		return piece
	}
	fmt.Println("Invalid input")
	return nil
}
