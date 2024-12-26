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
	successfulMoves := uint16(0)
	for {
		game.PrintBoard()
		player.printRack()
		fmt.Print("Enter a command (combine, insert, split, done): ")
		command, err := Reader.ReadString('\n')
		command = strings.TrimSpace(command)
		if err != nil || command == "done" {
			player.totalMoves = player.totalMoves + successfulMoves
			if successfulMoves == 0 {
				player.DealPiece(game.TakePiece())
			}
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

// func (p *player) insert(game Game) error {
// 	set, err := p.promptForSet(game)
// 	if err != nil {
// 		return err
// 	}
// 	piece, err := p.promptForPiece(nil)
// 	if err != nil {
// 		return err
// 	}
// 	if err := set.Insert(piece); err != nil {
// 		return err
// 	}
// 	p.removePiece(piece)
// 	return nil
// }

// func (player *player) combine(game Game) error {
// 	pieces := make([]Piece, 0)
// 	originalBoard := game.CloneBoard()
// 	for {
// 		fmt.Print("Select a piece <r#|s#,#|done> : ")
// 		input, err := Reader.ReadString('\n')
// 		if err != nil {
// 			return err
// 		}
// 		input = strings.TrimSpace(input)
// 		if input == "done" {
// 			break
// 		}
// 		piece, err := player.selectPiece(input, game)
// 		if err != nil {
// 			game.SetBoard(originalBoard)
// 			return err
// 		}
// 		pieces = append(pieces, piece)
// 		fmt.Printf("=== Selected Pieces ===\n%s=======================\n", (&set{tiles: pieces}).String())
// 	}
// 	set, err := Combine(pieces...)
// 	if err != nil {
// 		game.SetBoard(originalBoard)
// 		return err
// 	}
// 	player.removePiece(pieces...)
// 	game.AddSet(set)
// 	return nil
// }

// func (player *player) split(game Game) error {
// 	set, err := player.promptForSet(game)
// 	if err != nil {
// 		return err
// 	}
// 	piece, err := player.promptForPiece(nil)
// 	if err != nil {
// 		return err
// 	}
// 	splitSet, err := set.Split(piece)
// 	if err != nil {
// 		return err
// 	}
// 	game.AddSet(splitSet)
// 	player.removePiece(piece)
// 	return nil
// }

func parseCommand(input string, successfulMoves *uint16, player *player, game Game) {
	switch input {
	// case "combine":
	// 	if err := player.combine(game); err != nil {
	// 		fmt.Println(err)
	// 	} else {
	// 		*successfulMoves++
	// 	}
	// case "insert":
	// 	if err := player.insert(game); err != nil {
	// 		fmt.Println(err)
	// 	} else {
	// 		*successfulMoves++
	// 	}
	// case "split":
	// 	if err := player.split(game); err != nil {
	// 		fmt.Println(err)
	// 	} else {
	// 		*successfulMoves++
	// 	}
	case "help":
		fmt.Println("combine <r#|s#,#> <r#|s#,#> ... <r#|s#,#>")
		fmt.Println("insert <set> <piece>")
		fmt.Println("split <set> <piece>")
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
	if index < 0 || index >= len(p.rack) {
		return nil, errors.New(IndexOutOfBounds(len(p.rack)-1, "piece"))
	}
	return p.rack[index], nil
}

func (player *player) selectPiece(input string, game Game) (Piece, error) {
	input = strings.TrimSpace(input)
	start, selection := string(input[0]), string(input[1:])
	if start == "r" {
		return selectRackPiece(selection, player)
	}
	if start == "s" {
		return selectSetPiece(selection, game)
	}
	return nil, errors.New(InvalidPieceSelection)
}

func selectRackPiece(selection string, player *player) (Piece, error) {
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

func selectSetPiece(selection string, game Game) (Piece, error) {
	tuple := strings.Split(selection, ",")
	if len(tuple) != 2 {
		return nil, errors.New(InvalidNumberInput)
	}
	setIndex, err := parseInt(tuple[0])
	if err != nil {
		return nil, err
	}
	pieceIndex, err := parseInt(tuple[1])
	if err != nil {
		return nil, err
	}
	set, err := game.Set(setIndex)
	if err != nil {
		return nil, err
	}
	piece, err := set.Piece(pieceIndex)
	if err != nil {
		return nil, err
	}
	// if err := set.Remove(piece); err != nil {
	// 	return nil, err
	// }
	return piece, nil
}
