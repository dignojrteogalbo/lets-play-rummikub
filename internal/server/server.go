package server

import (
	"lets-play-rummikub/internal/model"
)

type Server struct {
	gameStarted   bool
	tilesShuffled bool
	tilesDealt    bool
	gameInstance  model.Game
	clients       map[*Client]model.Player
	receive       chan []byte
	register      chan *Client
	unregister    chan *Client
}

type ClientMessage struct {
	Client  *Client
	Message []byte
}

func NewServer(totalPlayers uint) *Server {
	return &Server{
		gameInstance: model.NewGame(totalPlayers),
		clients:      make(map[*Client]model.Player),
		receive:      make(chan []byte),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
	}
}

func (s *Server) Notify() {
	for client, player := range s.clients {
		gameState, err := s.gameInstance.MarshalJSON()
		if err == nil {
			client.send <- gameState
		}
		playerState, err := player.MarshalJSON()
		if err == nil {
			client.send <- playerState
		}
	}
}

func (s *Server) Run() {
	for {
		select {
		case client := <-s.register:
			s.clients[client] = s.gameInstance.Player(len(s.clients))
			currentBoard, err := s.gameInstance.MarshalJSON()
			if err == nil {
				client.send <- currentBoard
			}
		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
			}
		case message := <-s.receive:
			go s.handleCommands(message)
		}
	}
}

func (s *Server) handleCommands(message []byte) {
	switch string(message) {
	case "shuffle":
		if !s.tilesShuffled {
			s.gameInstance.Shuffle()
			s.tilesShuffled = true
		}
	case "deal":
		if !s.tilesDealt && len(s.clients) == s.gameInstance.TotalPlayers() {
			s.gameInstance.DealPieces()
			for client, player := range s.clients {
				playerState, err := player.MarshalJSON()
				if err == nil {
					client.send <- playerState
				}
			}
			s.tilesDealt = true
		} else {
			for client := range s.clients {
				client.send <- []byte("not enough players connected to deal pieces")
			}
		}
	case "start":
		if !s.gameStarted && len(s.clients) == s.gameInstance.TotalPlayers() {
			s.gameStarted = true
			s.gameInstance.Start(s)
		} else {
			for client := range s.clients {
				client.send <- []byte("not enough players to start game")
			}
		}
	default:
		for client := range s.clients {
			client.send <- message
		}
	}
}
