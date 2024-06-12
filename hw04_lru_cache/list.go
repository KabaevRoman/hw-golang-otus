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
	count int
}

func (l *ListImpl) Len() int {
	return l.count
}

func (l *ListImpl) Back() *ListItem {
	if l.list == nil {
		return nil
	}
	curr := l.list
	for curr.Next != nil {
		curr = curr.Next
	}
	return curr
}

func (l *ListImpl) Front() *ListItem {
	if l.list == nil {
		return nil
	}
	curr := l.list
	for curr.Prev != nil {
		curr = curr.Prev
	}
	return curr
}

func (l *ListImpl) PushFront(v interface{}) *ListItem {
	newNode := &ListItem{Value: v}
	l.count++
	if l.list == nil {
		l.list = newNode
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
		return newNode
	}
	back := l.Back()
	back.Next = newNode
	newNode.Prev = back
	return newNode
}

func (l *ListImpl) Remove(i *ListItem) {
	l.count--
	if i == l.list {
		l.list = l.list.Next
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
}

func NewList() List {
	return &ListImpl{}
}
