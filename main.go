package main

import (
	"fmt"
	"lets-play-rummikub/internal/model"
)

func main() {
	totalPlayers := 5
	game := model.NewGame(uint(totalPlayers))
	fmt.Printf("Started game with %d players\n", totalPlayers)
	game.Shuffle()
	fmt.Println("Shuffling tiles...")
	game.DealPieces()
	fmt.Println("Dealing tiles to players...")
	fmt.Println("Start game!")
	game.Start()
}
