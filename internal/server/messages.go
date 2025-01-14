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
	player, game, moveHistory := c.server.clients[c], c.server.game, c.server.history
	switch event.Command {
	case "name":
		command.SetName(player, event.Input).Invoke()
		c.send <- []byte(fmt.Sprintf(playerRenamed, player.Name()))
	case "combine":
		if playerCommand, err := command.Combine(player, game, event.Input); err == nil {
			playerCommand.Invoke()
			moveHistory.Push(playerCommand)
		} else {
			c.send <- []byte(fmt.Sprintf(commandError, event.Command, err.Error()))
		}
	case "insert":
		if playerCommand, err := command.Insert(player, game, event.Input); err == nil {
			playerCommand.Invoke()
			moveHistory.Push(playerCommand)
		} else {
			c.send <- []byte(fmt.Sprintf(commandError, event.Command, err.Error()))
		}
	case "remove":
		if playerCommand, err := command.Remove(game, event.Input); err == nil {
			playerCommand.Invoke()
			moveHistory.Push(playerCommand)
		} else {
			c.send <- []byte(fmt.Sprintf(commandError, event.Command, err.Error()))
		}
	case "split":
		if playerCommand, err := command.Split(game, event.Input); err == nil {
			playerCommand.Invoke()
			moveHistory.Push(playerCommand)
		} else {
			c.send <- []byte(fmt.Sprintf(commandError, event.Command, err.Error()))
		}
	default:
		c.send <- []byte(invalidCommand)
	}
}
