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

	Piece interface {
		IsJoker() bool
		IsSameColor(Piece) bool
		IsSameValue(Piece) bool
		IsSamePiece(Piece) bool
		Value() uint8
	}

	piece struct {
		value uint8
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
	if piece.value > 13 {
		return false
	}
	if piece.Color > ColorGreen {
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
		return p.Color == compare.(*piece).Color
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

func (p *piece) Value() uint8 {
	if !isValidPiece(p) {
		return 0
	}
	return p.value
}
