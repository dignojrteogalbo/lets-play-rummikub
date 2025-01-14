package main

import (
	"fmt"
	"lets-play-rummikub/internal/server"
	"net/http"
	"os"
)

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
	gameServer := server.NewServer(2)
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
