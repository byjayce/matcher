package internal

import "container/heap"

type Participant[A, B Identifier] struct {
	Self          A
	EngagedTo     *PreferredTarget[A, B]
	preferredList *list[B, A]
}

func NewParticipant[A, B Identifier](self A) *Participant[A, B] {
	return &Participant[A, B]{Self: self, preferredList: newList[B, A]()}
}

func (p *Participant[A, B]) Score(id string) (score int, isTop bool, ok bool) {
	target, ok := p.preferredList.set[id]
	if !ok {
		return 0, false, false
	}

	return target.Score, p.preferredList.list[0] == target, true
}

func (p *Participant[A, B]) Remove(id string) {
	for i, target := range p.preferredList.list {
		if target.Participant.Self.ID() == id {
			heap.Remove(p.preferredList, i)
			delete(p.preferredList.set, id)
		}
	}
}

func (p *Participant[A, B]) AddPreferredParticipant(score int, target *Participant[B, A]) {
	heap.Push(p.preferredList, &PreferredTarget[A, B]{Score: score, Participant: target})
}

func (p *Participant[A, B]) PopPreferredParticipant() *PreferredTarget[A, B] {
	target, ok := heap.Pop(p.preferredList).(*PreferredTarget[A, B])
	if !ok {
		return nil
	}
	return target
}

func (p *Participant[A, B]) Engage(score int, target *Participant[B, A]) *PreferredTarget[A, B] {
	before := p.EngagedTo
	p.EngagedTo = &PreferredTarget[A, B]{Score: score, Participant: target}
	p.Remove(target.Self.ID())
	return before
}
