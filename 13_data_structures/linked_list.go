package main

// Single Linked List
type LLNode[T any] struct {
	data T
	next *LLNode[T]
}

func NewLLNode[T any](data T) *LLNode[T] {
	return &LLNode[T]{
		data: data,
		next: nil,
	}
}

type LinkedList[T any] struct {
	head *LLNode[T]
	size int
}

func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{
		head: nil,
		size: 0,
	}
}

func (l *LinkedList[T]) Append(data T) {
	if l == nil {
		return
	}
	if l.head == nil {
		l.head = NewLLNode(data)
	} else {
		cur := l.head
		for cur.next != nil {
			cur = cur.next
		}
		cur.next = NewLLNode(data)
	}
	l.size++
}

func (l *LinkedList[T]) Delete(idx int) bool {
	// l might be null
	if l == nil {
		return false
	}
	// head might be null
	cur := l.head
	if cur == nil {
		return false
	}
	// idx might be out of range
	if idx < 0 || idx >= l.size {
		return false
	}
	if idx == 0 {
		l.head = l.head.next
	} else {
		for i := 0; i < idx-1; i++ {
			cur = cur.next
		}
		cur.next = cur.next.next
	}
	l.size--
	return true
}

func (l *LinkedList[T]) Update(idx int, newData T) bool {
	if l == nil {
		return false
	}
	if l.head == nil {
		return false
	}
	if idx < 0 || idx >= l.size {
		return false
	}
	cur := l.head
	for range idx {
		cur = cur.next
	}
	cur.data = newData
	return false
}

func (l *LinkedList[T]) Get(idx int) (T, bool) {
	var zero T
	if l == nil {
		return zero, false
	}
	if l.head == nil {
		return zero, false
	}
	if idx < 0 || idx >= l.size {
		return zero, false
	}
	cur := l.head
	for range idx {
		cur = cur.next
	}
	return cur.data, true
}

func (l *LinkedList[T]) Size() int {
	return l.size
}
