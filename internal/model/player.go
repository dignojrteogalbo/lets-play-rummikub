package model

type Player interface {}

func NewPlayer() Player {
	return new(Player)
}