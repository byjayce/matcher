package internal

type PreferredTarget[A, B Identifier] struct {
	Score       int
	Participant *Participant[B, A]
}
