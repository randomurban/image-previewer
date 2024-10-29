package cache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	length int
	front  *ListItem
	back   *ListItem
}

func (l *list) PushFront(v interface{}) *ListItem {
	n := new(ListItem)
	n.Value = v
	n.Next = l.front
	if l.length == 0 {
		l.back = n
	} else {
		l.front.Prev = n
	}
	l.front = n
	l.length++
	return n
}

func (l *list) PushBack(v interface{}) *ListItem {
	n := new(ListItem)
	n.Value = v
	n.Prev = l.back
	if l.length == 0 {
		l.front = n
	} else {
		l.back.Next = n
	}
	l.back = n
	l.length++
	return n
}

func (l *list) Disconnect(i *ListItem) {
	prev := i.Prev
	next := i.Next
	if prev == nil {
		l.front = next
	} else {
		prev.Next = next
	}
	if next == nil {
		l.back = prev
	} else {
		next.Prev = prev
	}
}

func (l *list) Remove(i *ListItem) {
	if l.length > 0 && i != nil {
		l.Disconnect(i)
		l.length--
	}
}

func (l *list) MoveToFront(i *ListItem) {
	if l.length > 1 && i != nil {
		l.Disconnect(i)
		i.Next = l.front
		l.front.Prev = i
		i.Prev = nil
		l.front = i
	}
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}
