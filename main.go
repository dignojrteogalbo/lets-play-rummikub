package main

import (
	"fmt"
	"lets-play-rummikub/internal/server"
)

func main() {
	server := &server.Server{Message: "Hello world!"}
	fmt.Printf(server.Message)
}
