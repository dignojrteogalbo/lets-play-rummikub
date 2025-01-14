package server

import (
	"lets-play-rummikub/internal/history"
	"lets-play-rummikub/internal/model"
)

type Server struct {
	gameStarted   bool
	tilesShuffled bool
	tilesDealt    bool
	game          model.Game
	history       history.Stack[history.Undoable]
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
		game:       model.NewGame(totalPlayers),
		clients:    make(map[*Client]model.Player),
		receive:    make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		history:    history.NewStack[history.Undoable](),
	}
}

func (s *Server) Notify(message ...string) {
	for client, player := range s.clients {
		gameState, err := s.game.MarshalJSON()
		if err == nil {
			client.send <- gameState
		}
		playerState, err := player.MarshalJSON()
		if err == nil {
			client.send <- playerState
		}
		if s.game.CurrentPlayer() == player {
			for _, m := range message {
				client.send <- []byte(m)
			}
		}
	}
}

func (s *Server) Run() {
	for {
		select {
		case client := <-s.register:
			s.clients[client] = s.game.Player(len(s.clients))
			currentBoard, err := s.game.MarshalJSON()
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
			s.game.Shuffle()
			s.tilesShuffled = true
		}
	case "deal":
		if !s.tilesDealt && len(s.clients) == s.game.TotalPlayers() {
			s.game.DealPieces()
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
		if !s.gameStarted && len(s.clients) == s.game.TotalPlayers() {
			s.gameStarted = true
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
