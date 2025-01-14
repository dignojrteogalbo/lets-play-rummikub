package server

import (
	"fmt"
	"lets-play-rummikub/internal/command"
)

type Event struct {
	Command string `json:"command"`
	Input   string `json:"input"`
}

const (
	commandError   = string("error performing %s: %s")
	playerRenamed  = string("your name has been set to: %s")
	invalidCommand = string("invalid command")
)

func (c *Client) handleCommand(event Event) {
	server, player, game, moveHistory := c.server, c.server.clients[c], c.server.game, c.server.history
	if server.gameStarted {
		switch event.Command {
		case "combine":
			if playerCommand, err := command.Combine(player, game, event.Input); err == nil {
				playerCommand.Invoke()
				moveHistory.Push(playerCommand)
			} else {
				c.send <- []byte(fmt.Sprintf(commandError, event.Command, err.Error()))
			}
			return
		case "insert":
			if playerCommand, err := command.Insert(player, game, event.Input); err == nil {
				playerCommand.Invoke()
				moveHistory.Push(playerCommand)
			} else {
				c.send <- []byte(fmt.Sprintf(commandError, event.Command, err.Error()))
			}
			return
		case "remove":
			if playerCommand, err := command.Remove(game, event.Input); err == nil {
				playerCommand.Invoke()
				moveHistory.Push(playerCommand)
			} else {
				c.send <- []byte(fmt.Sprintf(commandError, event.Command, err.Error()))
			}
			return
		case "split":
			if playerCommand, err := command.Split(game, event.Input); err == nil {
				playerCommand.Invoke()
				moveHistory.Push(playerCommand)
			} else {
				c.send <- []byte(fmt.Sprintf(commandError, event.Command, err.Error()))
			}
			return
		case "undo":
			command := moveHistory.Pop()
			if command != nil {
				command.Undo()
			}
			return
		}
	}
	switch event.Command {
	case "start":
		if !server.gameStarted && len(server.clients) == server.game.TotalPlayers() {
			server.gameStarted = true
			game.Notify(fmt.Sprintf("%s's turn\n", game.CurrentPlayer().Name()))
		} else {
			c.send <- []byte("not enough players to start game")
		}
	case "end":
		game.NextTurn()
	case "shuffle":
		if !server.tilesShuffled {
			game.Shuffle()
			server.tilesShuffled = true
			game.Notify()
		}
	case "deal":
		if !server.tilesDealt && len(server.clients) == server.game.TotalPlayers() {
			game.DealPieces()
			game.Notify()
			server.tilesDealt = true
		} else {
			c.send <- []byte("not enough players connected to deal pieces")
		}
	case "name":
		command.SetName(player, event.Input).Invoke()
		c.send <- []byte(fmt.Sprintf(playerRenamed, player.Name()))
	default:
		c.send <- []byte(invalidCommand)
	}
}
