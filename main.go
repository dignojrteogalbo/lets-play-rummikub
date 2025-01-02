package main

import (
	"fmt"
	"lets-play-rummikub/internal/model"
	"lets-play-rummikub/internal/server"
	"net/http"
	"os"
)

func playLocally() {
	totalPlayers := 5
	game := model.NewGame(uint(totalPlayers))
	fmt.Printf("Started game with %d players\n", totalPlayers)
	game.Shuffle()
	fmt.Println("Shuffling tiles...")
	game.DealPieces()
	fmt.Println("Dealing tiles to players...")
	fmt.Println("Start game!")
	game.Start(nil)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "template/home.html")
}

func main() {
	gameServer := server.NewServer(5)
	go gameServer.Run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWs(gameServer, w, r)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Print("ListenAndServe: ", err)
		os.Exit(0)
	}
}
