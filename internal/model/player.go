package model

import (
	"fmt"
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

func (p *player) StartTurn(game Game) {
	successfulMoves := uint16(0)
	for {
		p.printRack()
		input, err := Reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if err != nil || input == "done" {
			if successfulMoves == 0 {
				p.rack = append(p.rack, game.TakePiece())
			}
			p.totalMoves = p.totalMoves + successfulMoves
			break
		}
		command := strings.Split(input, " ")
		parseCommand(command, &successfulMoves, p, game)
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

func (p *player) removePiece(remove ...Piece) {
	for {
		if len(remove) == 0 {
			return
		}
		top := len(remove) - 1
		for index, tile := range p.rack {
			if tile == remove[top] {
				p.rack = append(p.rack[:index], p.rack[index+1:]...)
				remove = remove[:top]
				break
			}
		}
	}
}

// func playerCombine(player *player, game Game, options ...string) error {
// 	setsAndIndexesOrPieces := make([]any, 0)
// 	removePieces := make([]Piece, 0)
// 	for _, o := range options {
// 		start := string(o[0])
// 		if start == "r" {
// 			index, err := strconv.ParseInt(o[1:], 0, 16)
// 			if err != nil {
// 				return err
// 			}
// 			piece := player.selectPieceFromRack(int(index))
// 			setsAndIndexesOrPieces = append(setsAndIndexesOrPieces, piece)
// 			removePieces = append(removePieces, piece)
// 		} else if start == "s" {
// 			rest := o[1:]
// 			params := strings.Split(rest, ",")
// 			if len(params) < 2 {
// 				return errors.New(InvalidParameters)
// 			}
// 			setIndex, err := strconv.ParseInt(params[0], 0, 16)
// 			if err != nil {
// 				return err
// 			}
// 			set, err := game.BoardSet(int(setIndex))
// 			if err != nil {
// 				return err
// 			}
// 			pieceIndex, err := strconv.ParseInt(params[1], 0, 16)
// 			if err != nil {
// 				return err
// 			}
// 			setsAndIndexesOrPieces = append(setsAndIndexesOrPieces, set, pieceIndex)
// 		}
// 	}
// 	combineSet, err := Combine(setsAndIndexesOrPieces...)
// 	if err != nil {
// 		return err
// 	}
// 	game.AddSet(combineSet)
// 	player.removePiece(removePieces...)
// 	fmt.Println("combined into set")
// 	game.PrintBoard()
// 	return nil
// }

// func playerInsert(player *player, game Game, options ...string) error {
// 	if len(options) < 3 {
// 		return errors.New(InvalidParameters)
// 	}
// 	setIndex, err := strconv.ParseInt(options[0], 0, 16)
// 	if err != nil {
// 		return err
// 	}
// 	pieceIndex, err := strconv.ParseInt(options[1], 0, 16)
// 	if err != nil {
// 		return err
// 	}
// 	position, err := strconv.ParseInt(options[2], 0, 16)
// 	if err != nil {
// 		return err
// 	}
// 	piece := player.selectPieceFromRack(int(pieceIndex))
// 	insertSet, err := game.BoardSet(int(setIndex))
// 	if err != nil {
// 		return err
// 	}
// 	err = insertSet.Insert(piece, int(position))
// 	if err != nil {
// 		return err
// 	}
// 	player.removePiece(piece)
// 	fmt.Println("inserted into set")
// 	game.PrintBoard()
// 	return err
// }

// func playerRemove(player *player, game Game, options ...string) error {
// 	setIndex, err := strconv.ParseInt(options[0], 0, 16)
// 	if err != nil {
// 		return err
// 	}
// 	pieceIndex, err := strconv.ParseInt(options[1], 0, 16)
// 	if err != nil {
// 		return err
// 	}
// 	position, err := strconv.ParseInt(options[2], 0, 16)
// 	if err != nil {
// 		return err
// 	}
// 	pieces := make([]Piece, 0)
// 	for _, o := range options[3:] {
// 		index, err := strconv.ParseInt(o, 0, 16)
// 		if err != nil {
// 			return err
// 		}
// 		pieces = append(pieces, player.selectPieceFromRack(int(index)))
// 	}
// 	removeSet, err := game.BoardSet(int(setIndex))
// 	if err != nil {
// 		return err
// 	}
// 	set, err := removeSet.Remove(int(pieceIndex), int(position), pieces...)
// 	if err != nil {
// 		return err
// 	}
// 	game.AddSet(set)
// 	player.removePiece(pieces...)
// 	fmt.Println("removed from set")
// 	game.PrintBoard()
// 	return nil
// }

// func playerSplit(player *player, game Game, options ...string) error {
// 	if len(options) < 2 {
// 		return errors.New(InvalidParameters)
// 	}
// 	setIndex, err := strconv.ParseInt(options[0], 0, 16)
// 	if err != nil {
// 		return err
// 	}
// 	splitSet, err := game.BoardSet(int(setIndex))
// 	if err != nil {
// 		return err
// 	}
// 	pieceIndex, err := strconv.ParseInt(options[1], 0, 16)
// 	if err != nil {
// 		return err
// 	}
// 	piece := player.selectPieceFromRack(int(pieceIndex))
// 	set, err := splitSet.Split(piece)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	game.AddSet(set)
// 	player.removePiece(piece)
// 	fmt.Println("split set")
// 	game.PrintBoard()
// 	return nil
// }

func parseCommand(input []string, successfulMoves *uint16, player *player, game Game) {
	if len(input) < 2 {
		return
	}
	first, _ := input[0], input[1:]
	switch first {
	case "combine":
	// if err := playerCombine(player, game, options...); err == nil {
	// 	*successfulMoves++
	// } else {
	// 	fmt.Println(err)
	// }
	case "insert":
	// 	if err := playerInsert(player, game, options...); err == nil {
	// 		*successfulMoves++
	// 	} else {
	// 		fmt.Println(err)
	// 	}
	case "remove":
	// 	if err := playerRemove(player, game, options...); err == nil {
	// 		*successfulMoves++
	// 	} else {
	// 		fmt.Println(err)
	// 	}
	case "split":
	// 	if err := playerSplit(player, game, options...); err == nil {
	// 		*successfulMoves++
	// 	} else {
	// 		fmt.Println(err)
	// 	}
	case "help":
		fmt.Println("combine <r#|s#,#> <r#|s#,#> <r#|s#,#>")
		fmt.Println("insert <set> <rack> <position>")
		fmt.Println("remove <set> <piece> <position> <r1> <r2> ... <rn>")
		fmt.Println("split <set> <rack>")
	default:
		fmt.Println("invalid command")
	}
}

func (p *player) selectPieceFromRack(index int) Piece {
	if index < 0 || index >= len(p.rack) {
		return nil
	}
	return p.rack[index]
}
