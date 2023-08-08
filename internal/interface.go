package internal

type Identifier interface {
	// ID returns the unique ID of the participant.
	ID() string
}
