package server

import "lets-play-rummikub/internal/model"

type Server struct {
	players []model.Player
	Message string
}

func NewServer(message string, players []model.Player) *Server {
	return &Server{
		players: players,
		Message: message,
	}
}