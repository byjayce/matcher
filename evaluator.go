package matcher

import (
	"errors"
	"github.com/byjayce/matcher/internal"
)

var ErrImpossible = errors.New("impossible")

type Identifier = internal.Identifier

type Evaluator[T any] interface {
	Identifier
	// Evaluate returns the score of the other participant.
	// The score value is an integer greater than or equal to 0.
	// If error is returned and the error is ErrImpossible, The partner will not be in the preference list.
	// If error is returned and the error is not ErrImpossible, the partner will be in the preference list and the score will be 0.
	Evaluate(partner T) (int, error)
}
