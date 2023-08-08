package matcher

import (
	"errors"
	"github.com/changhoi/matcher/internal"
)

type Pair[Proposer, Recipient any] struct {
	Proposer  Proposer
	Recipient Recipient
}

type Failed[T any] struct {
	Participant T
	Reason      string
	CancelCount int
	RejectCount int
}

type MatchResult[Proposer, Recipient any] struct {
	Matched         []Pair[Proposer, Recipient]
	FailedProposer  []Proposer
	FailedRecipient []Recipient
}

func initParticipants[A Evaluator[B], B Evaluator[A]](proposers []A, recipients []B) ([]*internal.Participant[A, B], []*internal.Participant[B, A]) {
	var (
		proposerList  = make([]*internal.Participant[A, B], 0, len(proposers))
		recipientList = make([]*internal.Participant[B, A], 0, len(recipients))
	)

	for _, p := range proposers {
		proposerList = append(proposerList, internal.NewParticipant[A, B](p))
	}

	for _, r := range recipients {
		recipientList = append(recipientList, internal.NewParticipant[B, A](r))
	}

	return proposerList, recipientList
}

func initPreferredList[A Evaluator[B], B Evaluator[A]](proposers []*internal.Participant[A, B], recipients []*internal.Participant[B, A]) {
	for _, p := range proposers {
		for _, r := range recipients {
			rScore, rErr := p.Self.Evaluate(r.Self)
			switch {
			case errors.Is(rErr, ErrImpossible):
				continue
			default:
			}

			pScore, pErr := r.Self.Evaluate(p.Self)
			switch {
			case errors.Is(pErr, ErrImpossible):
				continue
			default:
			}

			p.AddPreferredParticipant(rScore, r)
			r.AddPreferredParticipant(pScore, p)
		}
	}
}

func initMap[A Evaluator[B], B Evaluator[A]](recipients []*internal.Participant[B, A]) map[string]*internal.Participant[B, A] {
	var (
		ret = make(map[string]*internal.Participant[B, A], len(recipients))
	)

	for _, r := range recipients {
		ret[r.Self.ID()] = r
	}

	return ret
}

func Match[Proposer Evaluator[Recipient], Recipient Evaluator[Proposer]](proposers []Proposer, recipients []Recipient) MatchResult[Proposer, Recipient] {
	proposerList, recipientList := initParticipants(proposers, recipients)
	initPreferredList(proposerList, recipientList)

	var (
		pairs               []Pair[Proposer, Recipient]
		notMatchedProposer  []Proposer
		notMatchedRecipient []Recipient
	)

	var (
		proposeQueue        = newMatchQueue(proposerList...)
		recipientCandidates = initMap(recipientList)
		stagedPairs         = make(map[string]Pair[Proposer, Recipient], len(recipientList))
	)

	// gale-shapley algorithm
	for {
		if proposeQueue.isEmpty() {
			break
		}

		p := proposeQueue.pop()

		for {
			preferred := p.PopPreferredParticipant()
			if preferred == nil {
				// 남은 선호자가 없음
				notMatchedProposer = append(notMatchedProposer, p.Self)
				break
			}

			pScore, isTop, ok := preferred.Participant.Score(p.Self.ID())
			if !ok {
				continue
			}

			if _, ok := recipientCandidates[preferred.Participant.Self.ID()]; ok {
				// 매칭된 적 없음
				preferred.Participant.Engage(pScore, p)
				p.Engage(preferred.Score, preferred.Participant)
				delete(recipientCandidates, preferred.Participant.Self.ID())
				if isTop {
					// 최종 확정
					pairs = append(pairs, Pair[Proposer, Recipient]{
						Proposer:  p.Self,
						Recipient: preferred.Participant.Self,
					})
					break
				}

				// 확정은 아님
				stagedPairs[preferred.Participant.Self.ID()] = Pair[Proposer, Recipient]{
					Proposer:  p.Self,
					Recipient: preferred.Participant.Self,
				}
				break
			}

			// 이미 매칭된 적 있음
			if _, ok := stagedPairs[preferred.Participant.Self.ID()]; ok {
				// 확정 전
				if pScore <= preferred.Participant.EngagedTo.Score {
					// 거절
					preferred.Participant.Remove(p.Self.ID())
					continue
				}

				// 상대 변경
				p.Engage(preferred.Score, preferred.Participant)
				before := preferred.Participant.Engage(pScore, p)

				if before != nil {
					// 다시 프로포즈
					proposeQueue.push(before.Participant)
					preferred.Participant.Remove(before.Participant.Self.ID())
				}

				if isTop {
					// 최종 확정
					pairs = append(pairs, Pair[Proposer, Recipient]{
						Proposer:  p.Self,
						Recipient: preferred.Participant.Self,
					})
					delete(stagedPairs, preferred.Participant.Self.ID())
					break
				}

				// 확정은 아님
				stagedPairs[preferred.Participant.Self.ID()] = Pair[Proposer, Recipient]{
					Proposer:  p.Self,
					Recipient: preferred.Participant.Self,
				}
			}

			// 이미 확정된 상대이므로 넘어감
		}
	}

	for _, r := range recipientCandidates {
		notMatchedRecipient = append(notMatchedRecipient, r.Self)
	}

	for _, p := range stagedPairs {
		// 마지막에 최종 확정
		pairs = append(pairs, p)
	}

	return MatchResult[Proposer, Recipient]{
		Matched:         pairs,
		FailedProposer:  notMatchedProposer,
		FailedRecipient: notMatchedRecipient,
	}
}
