package model

type Game interface {
	NextTurn()
	PlayMove()
	CheckPieces()
}

type gameInstance struct{}

func NewGame() Game {
	return new(gameInstance)
}

func (*gameInstance) NextTurn() {

}

func (*gameInstance) PlayMove() {

}

func (*gameInstance) CheckPieces() {

}
