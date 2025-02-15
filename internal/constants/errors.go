package constants

import "fmt"

const (
	InvalidSet              = string("set is invalid")
	InvalidPiece            = string("piece is invalid")
	InvalidCombineArguments = string("combine arguments must be pairs of type Set and int")
	InvalidPieceSelection   = string("invalid piece selection")
	InvalidSetSelection     = string("invalid set selection")
	InvalidNumberInput      = string("invalid input is not a number")
	InvalidBoard            = string("board is invalid")
	TooFewPieces            = string("not enough pieces to create set")
	TooFewArguments         = string("not enough arguments provided")
	CannotInsert            = string("piece cannot be inserted into set")
	CannotSplit             = string("set cannot be split")
	WrongColorForRun        = string("piece does not match the color of the run")
	WrongValueForGroup      = string("piece does not match the value of the group")
)

// Returns ("name" must be > "min" and < "max")
func IndexOutOfBounds(min, max int, name ...string) string {
	label := "index"
	if len(name) > 0 {
		label = name[0]
	}
	return fmt.Sprintf("%s must be > %d and < %d", label, min, max)
}
