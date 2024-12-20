package model

const (
	ValueJoker uint8 = iota
	ColorBlack Color = iota
	ColorBlue
	ColorRed
	ColorGreen
)

type (
	Color uint8
	Value uint8

	Game interface {
		NextTurn()
		PlayMove()
		CheckPieces()
	}

	gameInstance struct {
		tiles [106]Piece
	}

	Piece interface {
		IsJoker() bool
		IsSameColor(Piece) bool
		IsSameValue(Piece) bool
	}

	piece struct {
		Value uint8
		Color
	}
)

func NewPiece(v uint8, c Color) Piece {
	if v > 13 {
		return nil
	}
	if c > ColorGreen {
		return nil
	}
	return &piece{v, c}
}

func isValidPiece(p Piece) bool {
	if p == nil {
		return false
	}
	piece, ok := p.(*piece)
	if !ok {
		return false
	}
	if piece.Value > 13 {
		return false
	}
	if piece.Color > ColorGreen {
		return false
	}
	return true
}

func (p *piece) IsJoker() bool {
	if isValidPiece(p) {
		return p.Value == 0
	}
	return false
}

func (p *piece) IsSameColor(compare Piece) bool {
	if isValidPiece(p) && isValidPiece(compare) {
		return p.Color == compare.(*piece).Color
	}
	return false
}

func (p *piece) IsSameValue(compare Piece) bool {
	if isValidPiece(p) && isValidPiece(compare) {
		return p.Value == compare.(*piece).Value
	}
	return false
}

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
