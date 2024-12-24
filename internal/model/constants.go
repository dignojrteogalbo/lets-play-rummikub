package model

import "fmt"

const (
	InvalidSet              = string("set is invalid")
	InvalidPiece            = string("piece is invalid")
	InvalidCombineArguments = string("combine arguments must be pairs of type Set and int")
	InvalidPieceSelection   = string("invalid piece selection")
	InvalidSetSelection     = string("invalid set selection")
	InvalidNumberInput      = string("invalid input is not a number")
	TooFewPieces            = string("not enough pieces to create set")
	CannotInsert            = string("piece cannot be inserted into set")
	CannotSplit             = string("set cannot be split")
	WrongColorForRun        = string("piece does not match the color of the run")
	WrongValueForGroup      = string("piece does not match the value of the group")
)

func IndexOutOfBounds(max int, name ...string) string {
	label := "index"
	if len(name) > 0 {
		label = name[0]
	}
	return fmt.Sprintf("%s must be >= 0 and < %d", label, max)
}
