package model

type (
	Game interface {
		NextTurn()
		PlayMove()
		CheckPieces()
	}

	gameInstance struct {
		tiles [106]Piece
	}
)

func NewGame() Game {
	instance := new(gameInstance)
	index := 0
	for i := 0; i < 2; i++ {
		for color := ColorBlack; color <= ColorGreen; color++ {
			for value := uint8(1); value <= uint8(13); value++ {
				instance.tiles[index] = NewPiece(value, color)
				index++
			}
		}
		instance.tiles[index] = NewPiece(ValueJoker, ColorBlack)
		index++
	}
	return instance
}

func (*gameInstance) NextTurn() {

}

func (*gameInstance) PlayMove() {

}

func (*gameInstance) CheckPieces() {

}
