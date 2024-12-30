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
		rack []Piece
	}
)

func NewPlayer() Player {
	return &player{rack: make([]Piece, 0)}
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

func (player *player) StartTurn(game Game) {
	successfulMeld := uint16(0)
	for {
		game.PrintBoard()
		player.printRack()
		fmt.Print("Enter a command (combine, insert, split, done): ")
		command, err := Reader.ReadString('\n')
		command = strings.TrimSpace(command)
		if err != nil || command == "done" {
			if !game.IsValidBoard() {
				fmt.Println(InvalidBoard)
			} else {
				if successfulMeld == 0 {
					player.DealPiece(game.TakePiece())
				}
				break
			}
		}
		player.parseCommand(command, game, &successfulMeld)
	}
}

func (p *player) Score() uint16 {
	score := uint16(0)
	for _, piece := range p.rack {
		score = score + uint16(piece.Value())
	}
	return score
}

func (player *player) removePiece(pieces ...Piece) {
	for _, piece := range pieces {
		for index, p := range player.rack {
			if piece.IsSamePiece(p) {
				player.rack = append(player.rack[:index], player.rack[index+1:]...)
				break
			}
		}
	}
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

func (player *player) insert(game Game) error {
	set, err := player.promptForSet(game)
	if err != nil {
		return err
	}
	fmt.Print("Select a piece <r#|p#> : ")
	input, err := Reader.ReadString('\n')
	if err != nil {
		return err
	}
	input = strings.TrimSpace(input)
	piece, err := player.selectPiece(input, game)
	if err != nil {
		return err
	}
	fmt.Printf("Select index [0, %d] : ", set.Len())
	input, err = Reader.ReadString('\n')
	if err != nil {
		return err
	}
	input = strings.TrimSpace(input)
	index, err := parseInt(input)
	if err != nil {
		return err
	}
	insert, err := set.Insert(piece, index)
	if err != nil {
		return err
	}
	player.removePiece(piece)
	game.ReplaceSet(set, insert)
	return nil
}

func (player *player) combine(game Game) error {
	pieces := make([]Piece, 0)
	for {
		fmt.Print("Select a piece <r#|p#|done> : ")
		input, err := Reader.ReadString('\n')
		if err != nil {
			return err
		}
		input = strings.TrimSpace(input)
		if input == "done" {
			break
		}
		piece, err := player.selectPiece(input, game)
		if err != nil {
			return err
		}
		pieces = append(pieces, piece)
		fmt.Printf("=== Selected Pieces ===\n%s=======================\n", (&set{tiles: pieces}).String())
	}
	set := Combine(pieces...)
	player.removePiece(pieces...)
	game.RemovePieces(pieces...)
	game.AddSet(set)
	return nil
}

func (player *player) remove(game Game) error {
	set, err := player.promptForSet(game)
	if err != nil {
		return err
	}
	fmt.Printf("Select index [0, %d] : ", set.Len())
	input, err := Reader.ReadString('\n')
	if err != nil {
		return err
	}
	input = strings.TrimSpace(input)
	index, err := parseInt(input)
	if err != nil {
		return err
	}
	piece, err := set.Piece(index)
	if err != nil {
		return err
	}
	remove, err := set.Remove(piece)
	if err != nil {
		return err
	}
	game.AddLoosePiece(piece)
	game.ReplaceSet(set, remove)
	return nil
}

func (player *player) split(game Game) error {
	set, err := player.promptForSet(game)
	if err != nil {
		return err
	}
	if set.Len() < 2 {
		return errors.New(CannotSplit)
	}
	fmt.Printf("Select index [1, %d] : ", set.Len())
	input, err := Reader.ReadString('\n')
	if err != nil {
		return err
	}
	input = strings.TrimSpace(input)
	index, err := parseInt(input)
	if err != nil {
		return err
	}
	lowerSet, upperSet, err := set.Split(index)
	if err != nil {
		return err
	}
	game.ReplaceSet(set, lowerSet)
	game.AddSet(upperSet)
	return nil
}

func (player *player) parseCommand(input string, game Game, successfulMeld *uint16) {
	switch input {
	case "done":
	case "combine":
		if err := player.combine(game); err != nil {
			fmt.Println(err)
		} else if game.IsValidBoard() {
			*successfulMeld++
		}
	case "split":
		if err := player.split(game); err != nil {
			fmt.Println(err)
		}
	case "insert":
		if err := player.insert(game); err != nil {
			fmt.Println(err)
		} else if game.IsValidBoard() {
			*successfulMeld++
		}
	case "remove":
		if err := player.remove(game); err != nil {
			fmt.Println(err)
		}
	case "help":
		fmt.Println("valid commands are: combine, split, insert, remove, help, done")
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
	if index < 0 || index > len(p.rack)-1 {
		return nil, errors.New(IndexOutOfBounds(-1, len(p.rack), "piece"))
	}
	return p.rack[index], nil
}

func (player *player) selectPiece(input string, game Game) (Piece, error) {
	input = strings.TrimSpace(input)
	start, selection := string(input[0]), string(input[1:])
	switch start {
	default:
		return nil, errors.New(InvalidPieceSelection)
	case "r":
		return player.selectRackPiece(selection)
	case "p":
		return selectLoosePiece(selection, game)
	}
}

func (player *player) selectRackPiece(selection string) (Piece, error) {
	index, err := parseInt(selection)
	if err != nil {
		return nil, err
	}
	piece, err := player.Piece(index)
	if err != nil {
		return nil, err
	}
	return piece, nil
}

func selectLoosePiece(selection string, game Game) (Piece, error) {
	index, err := parseInt(selection)
	if err != nil {
		return nil, err
	}
	piece := game.Piece(index)
	if piece == nil {
		return nil, errors.New(InvalidPieceSelection)
	}
	return piece, nil
}
