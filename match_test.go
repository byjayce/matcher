package matcher

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
	"strconv"
	"testing"
)

type evaluator int

func (e evaluator) ID() string {
	return strconv.Itoa(int(e))
}

func (e evaluator) Evaluate(partner evaluator) (int, error) {
	v := int(partner + e)
	if v > 10 {
		return 0, ErrImpossible
	}

	return v, nil
}

func TestMatch(t *testing.T) {
	type args[Proposer Evaluator[Recipient], Recipient Evaluator[Proposer]] struct {
		proposers  []Proposer
		recipients []Recipient
	}
	type testCase[Proposer Evaluator[Recipient], Recipient Evaluator[Proposer]] struct {
		name string
		args args[Proposer, Recipient]
		want MatchResult[Proposer, Recipient]
	}
	tests := []testCase[evaluator, evaluator]{
		{
			name: "success to match",
			args: args[evaluator, evaluator]{
				proposers:  []evaluator{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
				recipients: []evaluator{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			want: MatchResult[evaluator, evaluator]{
				Matched: []Pair[evaluator, evaluator]{
					{Proposer: 1, Recipient: 9},
					{Proposer: 2, Recipient: 8},
					{Proposer: 3, Recipient: 7},
					{Proposer: 4, Recipient: 6},
					{Proposer: 5, Recipient: 5},
					{Proposer: 6, Recipient: 4},
					{Proposer: 7, Recipient: 3},
					{Proposer: 8, Recipient: 2},
					{Proposer: 9, Recipient: 1},
					{Proposer: 0, Recipient: 0},
				},
			},
		},
		{
			name: "fail to match",
			args: args[evaluator, evaluator]{
				proposers:  []evaluator{1},
				recipients: []evaluator{10},
			},
			want: MatchResult[evaluator, evaluator]{
				FailedProposer:  []evaluator{1},
				FailedRecipient: []evaluator{10},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Match(tt.args.proposers, tt.args.recipients); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Match() = %+v, want %+v\n cmp:-----\n %v", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}
