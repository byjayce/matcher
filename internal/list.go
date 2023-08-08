package internal

type list[A, B Identifier] struct {
	list []*PreferredTarget[A, B]
	set  map[string]*PreferredTarget[A, B]
}

func newList[A, B Identifier]() *list[A, B] {
	return &list[A, B]{set: make(map[string]*PreferredTarget[A, B])}
}

func (p *list[A, B]) Push(x any) {
	item := x.(*PreferredTarget[A, B])

	p.list = append(p.list, item)
	p.set[item.Participant.Self.ID()] = item
}

func (p *list[A, B]) Pop() any {
	n := len(p.list)
	if n == 0 {
		return nil
	}

	x := p.list[n-1]
	p.list = p.list[:n-1]
	delete(p.set, x.Participant.Self.ID())
	return x
}

func (p *list[A, B]) Len() int {
	return len(p.list)
}

func (p *list[A, B]) Less(i, j int) bool {
	n := len(p.list)
	if i >= n || j >= n {
		return false
	}

	return p.list[i].Score > p.list[j].Score
}

func (p *list[A, B]) Swap(i, j int) {
	n := len(p.list)
	if i >= n || j >= n {
		return
	}
	p.list[i], p.list[j] = p.list[j], p.list[i]
}
