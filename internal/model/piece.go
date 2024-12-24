package model

const (
	ValueJoker Value = iota
	ColorBlack Color = iota
	ColorBlue
	ColorRed
	ColorGreen
)

type (
	Color uint8
	Value uint8

	Piece interface {
		IsJoker() bool
		IsSameColor(Piece) bool
		IsSameValue(Piece) bool
		IsSamePiece(Piece) bool
		Value() Value
		Color() Color
	}

	piece struct {
		value Value
		color Color
	}
)

func NewPiece(v Value, c Color) Piece {
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
	if piece.value > 13 {
		return false
	}
	if piece.color > ColorGreen {
		return false
	}
	return true
}

func (p *piece) IsJoker() bool {
	if isValidPiece(p) {
		return p.value == 0
	}
	return false
}

func (p *piece) IsSameColor(compare Piece) bool {
	if isValidPiece(p) && isValidPiece(compare) {
		return p.color == compare.(*piece).color
	}
	return false
}

func (p *piece) IsSameValue(compare Piece) bool {
	if isValidPiece(p) && isValidPiece(compare) {
		return p.value == compare.(*piece).value
	}
	return false
}

func (p *piece) IsSamePiece(compare Piece) bool {
	return p == compare.(*piece)
}

func (p *piece) Value() Value {
	if !isValidPiece(p) {
		return 0
	}
	return p.value
}

func (p *piece) Color() Color {
	if !isValidPiece(p) {
		return 0
	}
	return p.color
}
