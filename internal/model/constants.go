package model

import "fmt"

const (
	InvalidSet              = string("set is invalid")
	InvalidPiece            = string("piece is invalid")
	InvalidCombineArguments = string("combine arguments must be pairs of type Set and int")
	InvalidParameters       = string("invalid command parameters, try \"help me\"")
	TooFewPieces            = string("not enough pieces to create set")
	CannotInsertIntoRun     = string("piece cannot be inserted into run")
	CannotSplitRun          = string("cannot split run")
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
