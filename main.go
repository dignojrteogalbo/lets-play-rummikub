package main

import (
	"fmt"
	"lets-play-rummikub/internal/model"
	"lets-play-rummikub/internal/server"
)

func main() {
	server := &server.Server{Message: "Hello world!\n"}
	fmt.Printf(server.Message)
	model.NewGame()
}
