package internal

type LinkedListNode[T comparable] struct {
	Value T
	Next  *LinkedListNode[T]
	Prev  *LinkedListNode[T]
}

func NewLinkedListNode[T comparable](value T) *LinkedListNode[T] {
	return &LinkedListNode[T]{Value: value}
}

type LinkedList[T comparable] struct {
	head *LinkedListNode[T]
	tail *LinkedListNode[T]
	len  int
}

func NewLinkedList[T comparable]() *LinkedList[T] {
	return &LinkedList[T]{}
}

func (l *LinkedList[T]) Len() int {
	return l.len
}

func (l *LinkedList[T]) First() *LinkedListNode[T] {
	return l.head
}

func (l *LinkedList[T]) Last() *LinkedListNode[T] {
	return l.tail
}

func (l *LinkedList[T]) AppendLast(node *LinkedListNode[T]) {
	if l.len == 0 {
		l.head = node
		l.tail = node
	} else {
		l.tail.Next = node
		node.Prev = l.tail
		l.tail = node
	}
	l.len++
}

func (l *LinkedList[T]) Remove(node *LinkedListNode[T]) {
	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		l.head = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		l.tail = node.Prev
	}
	node.Next = nil
	node.Prev = nil
	l.len--
}

func (l *LinkedList[T]) Find(value T) *LinkedListNode[T] {
	for node := l.head; node != nil; node = node.Next {
		if node.Value == value {
			return node
		}
	}
	return nil
}

func (l *LinkedList[T]) List() []T {
	var list []T
	for node := l.head; node != nil; node = node.Next {
		list = append(list, node.Value)
	}
	return list
}

func (l *LinkedList[T]) Clear() {
	l.head = nil
	l.tail = nil
	l.len = 0
}
