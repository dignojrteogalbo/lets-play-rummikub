package model

import (
	"bufio"
	"errors"
	"fmt"
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
		NextTurn()
		IsGameOver() bool
		Start()
		PrintBoard()
		TakePiece() Piece
		SetBoard(board []Set)
		CloneBoard() []Set
	}

	instance struct {
		tiles         []Piece
		sets          []Set
		players       []Player
		currentPlayer int
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
	instance.sets = make([]Set, 0)
	instance.createTiles()
	instance.createPlayers(int(totalPlayers))
	instance.currentPlayer = -1
	return instance
}

func (game *instance) CloneBoard() []Set {
	clone := make([]Set, len(game.sets))
	for i := range game.sets {
		clone[i] = game.sets[i].Clone()
	}
	return clone
}

func (game *instance) SetBoard(board []Set) {
	game.sets = board
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

func (g *instance) DealPieces() {
	for _, player := range g.players {
		for i := 0; i < 14; i++ {
			player.DealPiece(g.TakePiece())
		}
	}
}

func (g *instance) Set(index int) (Set, error) {
	if index < 0 || index > len(g.sets)-1 {
		return nil, errors.New(IndexOutOfBounds(-1, len(g.sets), "set"))
	}
	return g.sets[index], nil
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
	for i, s := range g.sets {
		fmt.Printf("--- [%d] ---\n%s", i, s.String())
	}
	fmt.Println("=== Board ===")
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
	g.sets = append(g.sets, set)
}

func (g *instance) Start() {
	for !g.IsGameOver() {
		g.NextTurn()
	}
	g.PrintScores()
}
