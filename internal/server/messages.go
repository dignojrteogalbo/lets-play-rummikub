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
	switch event.Command {
	case "combine":
		if game.CurrentPlayer() != player {
			return
		}
		if playerCommand, err := command.Combine(player, game, event.Input); err == nil {
			playerCommand.Invoke()
			moveHistory.Push(playerCommand)
		} else {
			c.send <- []byte(fmt.Sprintf(commandError, event.Command, err.Error()))
		}
	case "insert":
		if game.CurrentPlayer() != player {
			return
		}
		if playerCommand, err := command.Insert(player, game, event.Input); err == nil {
			playerCommand.Invoke()
			moveHistory.Push(playerCommand)
		} else {
			c.send <- []byte(fmt.Sprintf(commandError, event.Command, err.Error()))
		}
	case "remove":
		if game.CurrentPlayer() != player {
			return
		}
		if playerCommand, err := command.Remove(game, event.Input); err == nil {
			playerCommand.Invoke()
			moveHistory.Push(playerCommand)
		} else {
			c.send <- []byte(fmt.Sprintf(commandError, event.Command, err.Error()))
		}
	case "split":
		if game.CurrentPlayer() != player {
			return
		}
		if playerCommand, err := command.Split(game, event.Input); err == nil {
			playerCommand.Invoke()
			moveHistory.Push(playerCommand)
		} else {
			c.send <- []byte(fmt.Sprintf(commandError, event.Command, err.Error()))
		}
	case "undo":
		if game.CurrentPlayer() != player {
			return
		}
		command := moveHistory.Pop()
		if command != nil {
			command.Undo()
		}
	case "end":
		if game.CurrentPlayer() != player {
			return
		}
		if ok := game.NextTurn(); ok {
			moveHistory.Clear()
		}
	case "start":
		if !server.gameStarted && len(server.clients) == server.game.TotalPlayers() {
			server.gameStarted = true
			game.Notify(fmt.Sprintf("%s's turn\n", game.CurrentPlayer().Name()))
		} else {
			c.send <- []byte("not enough players to start game")
		}
	case "shuffle":
		if !server.tilesShuffled {
			game.Shuffle()
			server.tilesShuffled = true
			game.Notify()
		}
	case "deal":
		if !server.tilesDealt && len(server.clients) == server.game.TotalPlayers() {
			game.DealPieces()
			server.tilesDealt = true
			game.Notify()
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
