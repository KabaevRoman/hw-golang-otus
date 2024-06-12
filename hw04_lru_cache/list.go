package hw04lrucache

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

type ListImpl struct {
	list  *ListItem
	back  *ListItem
	count int
}

func (l *ListImpl) Len() int {
	return l.count
}

func (l *ListImpl) Back() *ListItem {
	return l.back
}

func (l *ListImpl) Front() *ListItem {
	return l.list
}

func (l *ListImpl) PushFront(v interface{}) *ListItem {
	newNode := &ListItem{Value: v}
	l.count++
	if l.list == nil {
		l.list = newNode
		l.back = newNode
		return newNode
	}
	oldFront := l.Front()
	l.list = newNode
	l.list.Next = oldFront
	oldFront.Prev = l.list
	return newNode
}

func (l *ListImpl) PushBack(v interface{}) *ListItem {
	newNode := &ListItem{Value: v}
	l.count++
	if l.list == nil {
		l.list = newNode
		l.back = newNode
		return newNode
	}
	back := l.back
	back.Next = newNode
	newNode.Prev = back
	l.back = newNode
	return newNode
}

func (l *ListImpl) Remove(i *ListItem) {
	l.count--
	if i == l.list {
		l.list = l.list.Next
	}
	if i == l.back {
		l.back = l.back.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
}

func (l *ListImpl) MoveToFront(i *ListItem) {
	if i == l.list {
		return
	}
	if i == l.back {
		l.back = l.back.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	i.Next = nil
	i.Prev = nil
	oldFront := l.Front()
	i.Next = oldFront
	oldFront.Prev = i
	l.list = i
}

func NewList() List {
	return &ListImpl{}
}
