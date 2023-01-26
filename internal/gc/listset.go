package gc

import "container/list"

type ListSet struct {
	list     *list.List
	elements map[string]*list.Element
}

func (l *ListSet) Add(value string) {
	if _, ok := l.elements[value]; !ok {
		l.elements[value] = l.list.PushBack(value)
	}
}

func (l *ListSet) Iter() *ListSetIterator {
	return &ListSetIterator{
		top:  l.list.Front(),
		list: l,
	}
}

type ListSetIterator struct {
	top  *list.Element
	list *ListSet
}

func (l *ListSetIterator) Scan() bool {
	return l.top != nil
}

func (l *ListSetIterator) Next() string {
	value := l.top.Value.(string)
	l.top = l.top.Next()

	return value
}

func (l *ListSetIterator) Add(value string) {
	l.list.Add(value)
}

func NewListset() *ListSet {
	return &ListSet{
		list:     list.New(),
		elements: make(map[string]*list.Element),
	}
}
