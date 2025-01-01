package model

import (
	"bufio"
	"errors"
	"fmt"
	"lets-play-rummikub/internal/history"
	"math/rand"
	"os"
	"time"
)

var (
	Reader = bufio.NewReader(os.Stdin)
)

type (
	Game interface {
		Shuffle()
		DealPieces()
		Set(index int) (Set, error)
		AddSet(Set)
		ReplaceSet(existing, replace Set)
		NextTurn()
		IsGameOver() bool
		Start()
		PrintBoard()
		TakePiece() Piece
		Piece(index int) Piece
		AddLoosePiece(piece Piece)
		RemovePieces(piece ...Piece)
		IsValidBoard() bool
		history.History
		history.Cloneable
	}

	instance struct {
		firstMeldComplete bool
		tiles         []Piece
		board         []Set
		loose         []Piece
		players       []Player
		currentPlayer int
		history.History
	}
)

func (g *instance) createTiles() {
	g.tiles = make([]Piece, 106)
	index := 0
	for i := 0; i < 2; i++ {
		for color := ColorBlack; color <= ColorGreen; color++ {
			for value := Value(1); value <= Value(13); value++ {
				g.tiles[index] = NewPiece(value, color)
				index++
			}
		}
		g.tiles[index] = NewPiece(ValueJoker, ColorBlack)
		index++
	}
}

func (g *instance) createPlayers(totalPlayers int) {
	g.players = make([]Player, 0, totalPlayers)
	for i := 0; i < int(totalPlayers); i++ {
		g.players = append(g.players, NewPlayer())
	}
}

func NewGame(totalPlayers uint) Game {
	if totalPlayers < 1 {
		return nil
	}
	instance := new(instance)
	instance.firstMeldComplete = false
	instance.board = make([]Set, 0)
	instance.createTiles()
	instance.createPlayers(int(totalPlayers))
	instance.currentPlayer = -1
	instance.History = history.NewHistory(instance)
	return instance
}

func (game *instance) Clone() history.Cloneable {
	newGame := new(instance)
	board := make([]Set, len(game.board))
	for i := range game.board {
		board[i] = game.board[i].Clone()
	}
	newGame.board = board
	loose := make([]Piece, len(game.loose))
	copy(loose, game.loose)
	newGame.loose = loose
	return newGame
}

func (game *instance) Restore(newGame history.Cloneable) {
	setInstance, ok := newGame.(*instance)
	if !ok || setInstance == nil {
		return
	}
	game.board = setInstance.board
	game.loose = setInstance.loose
}

func (game *instance) IsValidBoard() bool {
	for _, set := range game.board {
		if !set.IsValidSet() {
			return false
		}
	}
	return true
}

func (g *instance) Shuffle() {
	for i := range g.tiles {
		j := rand.Intn(i + 1)
		g.tiles[i], g.tiles[j] = g.tiles[j], g.tiles[i]
	}
}

func (g *instance) TakePiece() Piece {
	if len(g.tiles) == 0 {
		return nil
	}
	piece := g.tiles[len(g.tiles)-1]
	g.tiles = g.tiles[:len(g.tiles)-1]
	return piece
}

func (g *instance) Piece(index int) Piece {
	if len(g.loose) == 0 || index < 0 || index >= len(g.loose) {
		return nil
	}
	return g.loose[index]
}

func (g *instance) AddLoosePiece(piece Piece) {
	g.loose = append(g.loose, piece)
}

func (g *instance) RemovePieces(pieces ...Piece) {
	for _, piece := range pieces {
		for index, p := range g.loose {
			if piece.IsSamePiece(p) {
				g.loose = append(g.loose[:index], g.loose[index+1:]...)
				break
			}
		}
	}
}

func (g *instance) DealPieces() {
	for _, player := range g.players {
		for i := 0; i < 14; i++ {
			player.DealPiece(g.TakePiece())
		}
	}
}

func (g *instance) ReplaceSet(existing, replace Set) {
	for index, set := range g.board {
		if set == existing {
			g.board[index] = replace
			if replace.Len() == 0 {
				g.board = append(g.board[:index], g.board[index+1:]...)
			}
			return
		}
	}
}

func (g *instance) Set(index int) (Set, error) {
	if index < 0 || index > len(g.board)-1 {
		return nil, errors.New(IndexOutOfBounds(-1, len(g.board), "set"))
	}
	return g.board[index], nil
}

func (g *instance) NextTurn() {
	g.currentPlayer = (g.currentPlayer + 1) % len(g.players)
	fmt.Printf("Player #%d's turn\n", g.currentPlayer+1)
	g.players[g.currentPlayer].StartTurn(g)
}

func (g *instance) IsGameOver() bool {
	for _, p := range g.players {
		if p.Score() == 0 {
			return true
		}
	}
	return false
}

func (g *instance) PrintBoard() {
	fmt.Println("=== Board ===")
	for i, s := range g.board {
		fmt.Printf("--- [%d] ---\n%s", i, s.String())
	}
	fmt.Println("=== Board ===")
	fmt.Println("=== Loose Pieces ===")
	for i, p := range g.loose {
		fmt.Printf("[%d] %s\n", i, p.String())
	}
	fmt.Println("=== Loose Pieces ===")
}

func (g *instance) PrintScores() {
	fmt.Println("Scores:")
	for i, p := range g.players {
		go func() {
			score := p.Score()
			time.Sleep(time.Duration(time.Duration(p.Score()).Milliseconds()))
			if score == 0 {
				fmt.Printf("Winner is player #%d!\n\n", i+1)
			} else {
				fmt.Printf("Player #%d: %d\n", i+1, p.Score())
			}
		}()
	}
}

func (g *instance) AddSet(set Set) {
	g.board = append(g.board, set)
}

func (game *instance) restoreGameState() {
	currentPlayer := game.players[game.currentPlayer]
	for {
		moves, board := currentPlayer.Undo(), game.Undo()
		if board == nil && moves == nil {
			break
		}
	}
}

func (game *instance) hasSetOverThirty() bool {
	overThirtyPoints := false
	for _, set := range game.board {
		if set.Size() >= 30 {
			overThirtyPoints = true
		}
	}
	return overThirtyPoints
}

func (game *instance) hasSetWithJoker() bool {
	for _, set := range game.board {
		if set.NumberOfJokers() > 0 {
			return true
		}
	}
	return false
}

func (g *instance) Start() {
	for !g.IsGameOver() {
		g.NextTurn()
		if !g.IsValidBoard() {
			fmt.Println("board has invalid sets")
			g.restoreGameState()
		} else if !g.firstMeldComplete {
			if g.hasSetWithJoker() {
				fmt.Println("initial meld cannot contain joker")
				g.restoreGameState()
			} else if !g.hasSetOverThirty() {
				fmt.Println("initial meld must sum > 30")
				g.restoreGameState()
			} else {
				g.firstMeldComplete = true
			}
		}
		g.Clear()
		g.players[g.currentPlayer].Clear()
	}
	g.PrintScores()
}
