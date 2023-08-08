package matcher

import "container/list"

type matchQueue[T any] struct {
	list *list.List
}

func newMatchQueue[T any](data ...T) *matchQueue[T] {

	ret := &matchQueue[T]{
		list: list.New(),
	}

	for _, v := range data {
		ret.push(v)
	}

	return ret
}

func (q *matchQueue[T]) isEmpty() bool {
	return q.list.Len() == 0
}

func (q *matchQueue[T]) push(v T) {
	q.list.PushBack(v)
}

func (q *matchQueue[T]) pop() T {
	e := q.list.Front()
	q.list.Remove(e)
	return e.Value.(T)
}
