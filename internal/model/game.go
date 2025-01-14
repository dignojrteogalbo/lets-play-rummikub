package model

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"lets-play-rummikub/internal/constants"
	"lets-play-rummikub/internal/event"
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
		Start(event.Listener, history.Stack[history.Undoable])
		PrintBoard()
		TakePiece() Piece
		HasPiece
		AddLoosePiece(piece Piece)
		RemovePieces(piece ...Piece)
		IsValidBoard() bool
		CurrentPlayer() Player
		Player(index int) Player
		TotalPlayers() int
		MarshalJSON() ([]byte, error)
		Notify(message ...string)
		Clone() Game
		Restore(game Game)
	}

	instance struct {
		event.Listener
		firstMeldComplete bool
		tiles             []Piece
		board             []Set
		loose             []Piece
		players           []Player
		currentPlayer     int
	}
)

func (g *instance) Notify(messages ...string) {
	if g.Listener != nil {
		g.Listener.Notify(messages...)
	}
}

func (g *instance) MarshalJSON() ([]byte, error) {
	output := struct {
		Board  []Set   `json:"board"`
		Pieces []Piece `json:"piece"`
	}{
		g.board,
		g.loose,
	}
	return json.Marshal(output)
}

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
		player := NewPlayer()
		player.SetName(fmt.Sprintf("Player %d", i+1))
		g.players = append(g.players, player)
	}
}

func NewGame(totalPlayers uint) Game {
	if totalPlayers < 1 {
		return nil
	}
	instance := new(instance)
	instance.firstMeldComplete = false
	instance.board = make([]Set, 0)
	instance.loose = make([]Piece, 0)
	instance.createTiles()
	instance.createPlayers(int(totalPlayers))
	instance.currentPlayer = 0
	return instance
}

func (game *instance) Clone() Game {
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

func (game *instance) Restore(restore Game) {
	restoreGame, ok := restore.(*instance)
	if !ok || restoreGame == nil {
		return
	}
	game.board = restoreGame.board
	game.loose = restoreGame.loose
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

func (g *instance) Piece(index int) (Piece, error) {
	if index < 0 || index >= len(g.loose) {
		return nil, errors.New(constants.InvalidPieceSelection)
	}
	return g.loose[index], nil
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
		return nil, errors.New(constants.InvalidSetSelection)
	}
	return g.board[index], nil
}

func (g *instance) Player(index int) Player {
	if index < 0 || index >= len(g.players) {
		return nil
	}
	return g.players[index]
}

func (g *instance) CurrentPlayer() Player {
	return g.players[g.currentPlayer]
}

func (g *instance) TotalPlayers() int {
	return len(g.players)
}

func undoMoves(history history.Stack[history.Undoable], game Game) {
	for {
		if game.IsValidBoard() {
			return
		}
		command := history.Pop()
		if command == nil {
			return
		}
		command.Undo()
	}
}

func (g *instance) NextTurn(history history.Stack[history.Undoable]) {
	fmt.Printf("Player #%d's turn\n", g.currentPlayer+1)
	currentPlayer := g.CurrentPlayer()
	playerScore := currentPlayer.Score()
	g.CurrentPlayer().StartTurn(g)
	if !g.IsValidBoard() {
		fmt.Println("board has invalid sets")
		undoMoves(history, g)
	} else if !g.firstMeldComplete {
		if g.hasSetWithJoker() {
			fmt.Println("initial meld cannot contain joker")
			undoMoves(history, g)
		} else if !g.hasSetOverThirty() {
			fmt.Println("initial meld must sum > 30")
			undoMoves(history, g)
		} else {
			g.firstMeldComplete = true
		}
	}
	if playerScore >= currentPlayer.Score() {
		currentPlayer.DealPiece(g.TakePiece())
	}
	g.currentPlayer = (g.currentPlayer + 1) % len(g.players)
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

func (g *instance) Start(listener event.Listener, history history.Stack[history.Undoable]) {
	if g.Listener == nil && listener != nil {
		g.Listener = listener
	}
	for !g.IsGameOver() {
		g.NextTurn(history)
	}
	g.PrintScores()
}
